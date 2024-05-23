package pkg

import (
	"fmt"
	"log"
	"net"
)

type EchoServer struct{}

func (s *EchoServer) Start() {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go readWriteLoop(conn)
	}
}

func readWriteLoop(conn net.Conn) {
	for {
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		if err != nil {
			log.Fatal(err)
		}

		data := string(buf[:n])

		log.Printf("Server received '%v'", data)

		respData := []byte(fmt.Sprintf("You sent me: %v", data))

		_, err = conn.Write(respData)

		if err != nil {
			log.Fatal(err)
		}
	}
}