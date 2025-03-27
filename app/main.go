package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	tcpAddr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 4221,
	}
	l, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Printf("Listening on port: %d...\n", tcpAddr.Port)
	for {
		tcpConn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(c net.TCPConn) {
			fmt.Println("Incoming request...")
			statusLine, err := bufio.NewReader(&c).ReadString('\n')
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
				c.Close()
				os.Exit(1)
			}
			fmt.Printf("Status: %s\n", statusLine)
			parts := strings.Split(statusLine, " ")
			if parts[1] == "/" {
				c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			} else {
				c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			c.Close()
		}(*tcpConn)
	}
}
