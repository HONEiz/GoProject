package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to server!")
		return
	} else {
		fmt.Println("Connected to server!")
	}

	connReader := bufio.NewReader(conn)
	localReader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan nama anda : ")
	nama, err := localReader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the message!")
	} else {
		fmt.Println("Nama yang anda pilih :", strings.TrimSpace(nama))
	}

	_, errWrite := conn.Write([]byte(nama))
	nama = strings.TrimSpace(nama)
	if errWrite != nil {
		fmt.Println("Error dalam mengirim")
		return
	} else {
		fmt.Println("Nama " + nama + " terkirim")
	}

	go func() {
		for {
			message, err := connReader.ReadString('\n')
			if err != nil {
				fmt.Print("Terputus dari server.")
				return
			}
			fmt.Println("")
			fmt.Print(message)
		}
	}()

	for {
		fmt.Print("Masukkan pesan: ")
		message, _ := localReader.ReadString('\n')
		// fix loop kirim pesan

		message = strings.TrimSpace(message) // Hapus whitespace di awal/akhir pesan

		if message == "" {

			continue // Kembali ke awal loop
		} else {
			_, err = conn.Write([]byte(message + "\n"))
		}
		if err != nil {
			fmt.Println("")
			return
		}
	}
}
