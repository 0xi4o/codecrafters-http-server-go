package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const CRLF = "\r\n"

func main() {
	directory := flag.String("directory", "", "Path of the directory to look for files")
	flag.Parse()

	tcpAddr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 4221,
	}
	listener, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	flags := map[string]string{"directory": *directory}

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
			response := SerializeResponse(request, flags)
			c.Write([]byte(response))
			c.Close()
		}(*tcpConn)
	}
}
