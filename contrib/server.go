package main

import (
	"encoding/binary"
	"io"
	"fmt"
	"net"
)

func handleRequests(conn net.Conn) {
	hdr := make([]byte, 2)
	for {
		l, err := conn.Read(hdr)
		if err == io.EOF {
			conn.Close()
			return
		}
		if l != 2 || err != nil {
			panic(fmt.Sprintf("S: %v %v", l, err))
		}
		b := make([]byte, binary.BigEndian.Uint16(hdr))
		l, err = conn.Read(b)
		if l != len(b) || err != nil {
			panic(fmt.Sprintf("S: %v %v", l, err))
		}
		conn.Write(append(hdr, b...))
	}
}

func main() {
	serverSock, err := net.Listen("tcp", "localhost:7999")
	if err != nil {
		panic(fmt.Sprintf("S: %v", err))
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
