package main

import (
	"github.com/gliderlabs/ssh"
	"log"
	"net/http"
)

func main() {
	tunnel := make(chan Tunnel)

	go func() {
		handler := http.HandlerFunc(makeApiHandler(tunnel))
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(ssh.ListenAndServe(":2222", makeSSHHandler(tunnel)))
}
