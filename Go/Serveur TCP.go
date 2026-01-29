package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"project/vision"
	"runtime"
	"sync"
	"time"
)

const (
	addr          = ":9000"
	maxImageBytes = 20 << 20 // 20 MB
)

// ======== JOB / RESULT ========

type Job struct {
	Payload []byte
	ReplyCh chan Result
}

type Result struct {
	Data []byte
	Err  error
}

// ======== WORKER POOL ========

func startWorkerPool(n int) chan<- Job {
	jobs := make(chan Job, 2*n)

	for i := 0; i < n; i++ {
		go func(id int) {
			for job := range jobs {
				out, err := vision.ProcessImage(job.Payload)
				job.ReplyCh <- Result{Data: out, Err: err}
			}
		}(i)
	}
	return jobs
}

// ======== TCP FRAMING (length-prefix) ========

func readFrame(r *bufio.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}

	n := binary.BigEndian.Uint32(lenBuf[:])
	if n == 0 || n > maxImageBytes {
		return nil, fmt.Errorf("invalid frame size: %d", n)
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func writeFrame(w io.Writer, payload []byte) error {
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(payload)))

	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

// ======== CLIENT HANDLER ========

func handleClient(conn net.Conn, jobs chan<- Job, id string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		payload, err := readFrame(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("[%s] disconnected", id)
			} else {
				log.Printf("[%s] read error: %v", id, err)
			}
			return
		}

		replyCh := make(chan Result, 1)
		jobs <- Job{Payload: payload, ReplyCh: replyCh}

		res := <-replyCh
		if res.Err != nil {
			log.Printf("[%s] processing error: %v", id, res.Err)
			_ = writeFrame(conn, []byte("ERROR"))
			continue
		}

		_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := writeFrame(conn, res.Data); err != nil {
			log.Printf("[%s] write error: %v", id, err)
			return
		}
	}
}

// ======== MAIN ========

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	workers := runtime.NumCPU()

	jobs := startWorkerPool(workers)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer ln.Close()

	log.Printf("TCP server listening on %s (workers=%d)", addr, workers)

	var mu sync.Mutex
	connID := 0

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		mu.Lock()
		connID++
		id := fmt.Sprintf("client-%d-%s", connID, conn.RemoteAddr())
		mu.Unlock()

		log.Printf("[%s] connected", id)
		go handleClient(conn, jobs, id)
	}
}
