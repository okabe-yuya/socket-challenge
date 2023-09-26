package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

const (
	SERVER_HOST = "127.0.0.1"
	SERVER_PORT = 8001
)

func main() {
	fd, err := listenSocket(SERVER_HOST, SERVER_PORT)
	if err != nil {
		fmt.Println("[ERROR] Caused error:", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	fmt.Printf("[INFO] Server is listening on %s:%d...\n", SERVER_HOST, SERVER_PORT)

	for {
		clientFd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("[ERROR] Failed accept connection: ", err)
			continue
		}

		fmt.Println("[INFO] Accepted client")
		go handleClient(clientFd)
	}
}

func listenSocket(host string, port int) (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)
	if err != nil {
		return 0, errors.New("Failed create socket")
	}

	addr := syscall.SockaddrInet4{ Port: port }
	copy(addr.Addr[:], net.ParseIP(host).To4())
	if err := syscall.Bind(fd, &addr); err != nil {
		return 0, errors.New("Failed binding socket")
	}

	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return 0, errors.New("Failed listen socket")
	}
	return fd, nil
}

func handleClient(clientFd int) {
	buffer := make([]byte, 1024)
	soundFile, err := os.Open("./static/birdland1.mp3")
	if err != nil {
		fmt.Println("[ERROR] Failed Open file: ", err)
		os.Exit(1)
	}
	defer soundFile.Close()

	for {
		n, err := soundFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("[ERROR] Failed Read buffer: ", err)
			return
		}
		if _, err := syscall.Write(clientFd, buffer[:n]); err != nil {
			fmt.Println("[ERROR] Failed Write buffer to file descriptor: ", err)
			return
		}
	}

	if _, err := syscall.Write(clientFd, []byte("EOF")); err != nil {
		fmt.Println("[ERROR] Failed Write 'EOF' to file descriptor: ", err)
		os.Exit(1)
	}
	syscall.Close(clientFd)
}
