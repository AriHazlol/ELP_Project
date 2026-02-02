package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net"
	"runtime"
)

const addr = ":9000"

type Job struct {
	Payload []byte
	ReplyCh chan []byte
}

// --- LOGIQUE DE DESSIN IDENTIQUE AU MAIN-2 ---

func drawLinesCustom(img *image.RGBA, lines []Line) {
	if len(lines) == 0 {
		return
	}

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	halfHeight := height / 2

	var leftLines, rightLines []Line

	// 1. Séparation des lignes par angle
	for _, line := range lines {
		if line.Theta < 1.5 {
			leftLines = append(leftLines, line)
		} else if line.Theta > 1.6 {
			rightLines = append(rightLines, line)
		}
	}

	// 2. Calcul de la moyenne et dessin
	renderAverage := func(group []Line, col color.RGBA) {
		if len(group) == 0 {
			return
		}

		var sRho, sTheta float64
		for _, l := range group {
			sRho += l.Rho
			sTheta += l.Theta
		}
		avgR := sRho / float64(len(group))
		avgT := sTheta / float64(len(group))

		cosT, sinT := math.Cos(avgT), math.Sin(avgT)

		for x := 0; x < width; x++ {
			y := (avgR - float64(x)*cosT) / sinT
			if int(y) >= halfHeight && int(y) < height {
				for dy := -1; dy <= 1; dy++ { // Épaisseur
					if int(y)+dy < height && int(y)+dy >= 0 {
						img.SetRGBA(x, int(y)+dy, col)
					}
				}
			}
		}
	}

	renderAverage(leftLines, blue)
	renderAverage(rightLines, red)
}

// Helper pour convertir RGBA vers Gray
func convertToGrayStruct(img *image.RGBA) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			gray.SetGray(x, y, color.Gray{Y: c.R})
		}
	}
	return gray
}

// --- PIPELINE DE TRAITEMENT ---

func ProcessImage(imgData []byte) ([]byte, error) {
	// 1. Décodage
	imgSrc, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	bounds := imgSrc.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, imgSrc, bounds.Min, draw.Src)

	numGoroutines := runtime.NumCPU()

	// ÉTAPES DU PIPELINE
	ConcurrentGrayscale(rgbaImg, numGoroutines)
	applyFixedTriangleROI(rgbaImg, numGoroutines)

	grayImg := convertToGrayStruct(rgbaImg)
	blurredImg := GaussianBlurParallel(grayImg)
	edgeImg := CannyEdgeDetection(blurredImg, 50.0, numGoroutines)
	lines := HoughLines(edgeImg, 1, 180, 100, 10)

	// RÉSULTAT FINAL (Dessin sur l'original couleur)
	finalImg := image.NewRGBA(bounds)
	draw.Draw(finalImg, bounds, imgSrc, bounds.Min, draw.Src)

	drawLinesCustom(finalImg, lines) // Utilisation de la logique de moyenne

	// Encodage
	var outBuf bytes.Buffer
	if format == "png" {
		png.Encode(&outBuf, finalImg)
	} else {
		jpeg.Encode(&outBuf, finalImg, nil)
	}
	return outBuf.Bytes(), nil
}

// --- PARTIE RÉSEAU ---

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	jobs := make(chan Job, 10)

	// Workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for job := range jobs {
				res, err := ProcessImage(job.Payload)
				if err != nil {
					fmt.Println("Erreur:", err)
					job.ReplyCh <- nil
				} else {
					job.ReplyCh <- res
				}
			}
		}()
	}

	ln, _ := net.Listen("tcp", addr)
	fmt.Println("Serveur prêt sur le port 9000")

	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			defer c.Close()
			reader := bufio.NewReader(c)
			for {
				var lenBuf [4]byte
				if _, err := io.ReadFull(reader, lenBuf[:]); err != nil { return }
				length := binary.BigEndian.Uint32(lenBuf[:])

				data := make([]byte, length)
				if _, err := io.ReadFull(reader, data); err != nil { return }

				replyCh := make(chan []byte, 1)
				jobs <- Job{Payload: data, ReplyCh: replyCh}
				
				result := <-replyCh
				if result != nil {
					var respLen [4]byte
					binary.BigEndian.PutUint32(respLen[:], uint32(len(result)))
					c.Write(respLen[:])
					c.Write(result)
				}
			}
		}(conn)
	}
}
