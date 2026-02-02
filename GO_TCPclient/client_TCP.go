package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Erreur: Le serveur n'est pas lancé !")
		return
	}
	defer conn.Close()

	imgData, err := os.ReadFile("road.png") 
	if err != nil {
		fmt.Println("Erreur: Impossible de trouver road.png dans ce dossier.")
		return
	}

	sizeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuf, uint32(len(imgData)))
	conn.Write(sizeBuf)

	conn.Write(imgData)
	fmt.Println("Image PNG envoyée au serveur...")

	io.ReadFull(conn, sizeBuf)
	replySize := binary.BigEndian.Uint32(sizeBuf)
	replyData := make([]byte, replySize)
	io.ReadFull(conn, replyData)

	os.WriteFile("resultat_serveur.png", replyData, 0644)
	fmt.Println("Succès")
}
