package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mutex sync.Mutex

func main() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen!")
		os.Exit(1)
	} else {
		fmt.Println("Listening to port 9090")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to accept connection!")
			os.Exit(1)
		} else {
			go handleClient(conn)
		}
	}
}

// Handle koneksi client
func handleClient(conn net.Conn) {
	// Membaca nama client
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error membaca nama client:", err)
		return
	}

	// Menyimpan nama client
	name = strings.TrimSpace(name)
	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	fmt.Printf("Client %s terhubung.\n", name)

	// Menyebarkan pesan ke semua client
	go broadcast(fmt.Sprintf("%s bergabung ke chat\n", name), conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s terputus.\n", name)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			broadcast(fmt.Sprintf("%s keluar dari chat\n", name), conn)
			return
		}

		message = strings.TrimSpace(message)
		if message != "" {
			broadcast(fmt.Sprintf("%s: %s\n", name, message), conn)
		}
	}
}

// Mengirim pesan ke semua client kecuali pengirim
func broadcast(message string, sender net.Conn) {
	mutex.Lock()
	for conn := range clients {
		if conn != sender {
			fmt.Fprintf(conn, "%s", message)
		}
	}
	mutex.Unlock()
}

// Sumber referensi:
// https://www.kelche.co/blog/go/socket-programming/
