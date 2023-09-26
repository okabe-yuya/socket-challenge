package main

import (
	"fmt"
	"os"
	"syscall"
	"net"
)

func main() {
	// TCPソケットを作成
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// ソケットにアドレスを割り当て(bind)
	addr := syscall.SockaddrInet4{ Port: 8001 }
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	if err := syscall.Bind(fd, &addr); err != nil {
		fmt.Println("Error binding socket:", err)
		os.Exit(1)
	}

	// ソケットのリッスンを開始
	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		fmt.Println("Error listen socket:", err)
		os.Exit(1)
	}

	// 無限ループを利用してクライアントの受付を開始
	fmt.Println("Server is listening on 0.0.0.0:8081")
	for {
		clientFd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}

		// クライアント受付後、並列処理でレスポンスを書き込み
		go func(fd int) {
			defer syscall.Close(fd)

			data := []byte("hello world")
			if _, err := syscall.Write(fd, data); err != nil {
				fmt.Println("Error writing: ", err)
				return
			}
		}(clientFd)
	}
}
