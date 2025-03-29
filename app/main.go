package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

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
			r := make([]byte, 1024)
			_, err := bufio.NewReader(&c).Read(r)
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
				c.Close()
				os.Exit(1)
			}
			request := DeserializeRequest(string(r))
			pathParts := strings.Split(request.Path, "/")
			switch {
			case pathParts[1] == "echo":
				response := SerializeResponse(request)
				c.Write([]byte(fmt.Sprintf("%s %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.Body)))
			case request.Path == "/":
				response := SerializeResponse(request)
				c.Write([]byte(fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)))
			default:
				response := SerializeResponse(request)
				c.Write([]byte(fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)))
			}
			c.Close()
		}(*tcpConn)
	}
}
