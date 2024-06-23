package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s1 := server(":3333")
	_ = server(":3334", ":3333")
	_ = server(":3335", ":3333", ":3334")

	time.Sleep(500 * time.Millisecond)

	if err := s1.broadcast([]byte("Hello, World!")); err != nil {
		log.Fatalf("Failed to broadcast: %v\n", err)
	}

	exitCh := make(chan os.Signal, 1)

	// registers the given channel to receive os/unix notifications
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-exitCh

	log.Println("Shutting down...")
}

func server(addr string, connectAddresses ...string) *Server {
	tcpOpts := TCPTransportOpts{
		listenAddress: addr,
		handshake:     defaultHandshake,
	}
	opts := ServerOpts{
		Transport:        NewTCPTransport(tcpOpts),
		ConnectAddresses: connectAddresses,
	}
	s := NewServer(opts)

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	return s
}
