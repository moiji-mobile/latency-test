package main

import (
	"io"
	"fmt"
	"net"
)

func handleRequests(conn net.Conn) {
	b := make([]byte, 16)
	for {
		l, err := conn.Read(b)
		if err == io.EOF {
			conn.Close()
			return
		}
		if l != len(b) || err != nil {
			panic(fmt.Sprintf("%v %v", l, err))
		}
		conn.Write(b)
	}
}

func main() {
	serverSock, err := net.Listen("tcp", "localhost:7999")
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	defer serverSock.Close()

	for {
		clientSock, err := serverSock.Accept()
		if err != nil {
			continue
		}
		go handleRequests(clientSock)
	}
}
