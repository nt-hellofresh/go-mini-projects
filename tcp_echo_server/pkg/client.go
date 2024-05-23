package pkg

import (
	"log"
	"net"
)

func RunClient() {
	conn, err := net.Dial("tcp", ":3000")

	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("Foo"))

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if err != nil {
		log.Fatal(err)
	}

	response := string(buf[:n])

	log.Printf("Response from server: '%v'", response)
}
