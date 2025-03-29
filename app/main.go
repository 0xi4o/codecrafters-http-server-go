package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const CRLF = "\r\n"

func main() {
	tcpAddr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 4221,
	}
	listener, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Printf("Listening on port: %d...\n", tcpAddr.Port)
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(c net.TCPConn) {
			fmt.Println("Incoming request...")
			r := make([]byte, 1024)
			_, err := bufio.NewReader(&c).Read(r)
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
				c.Close()
				os.Exit(1)
			}
			request := DeserializeRequest(string(r))
			response := SerializeResponse(request)
			c.Write([]byte(response))
			c.Close()
		}(*tcpConn)
	}
}
