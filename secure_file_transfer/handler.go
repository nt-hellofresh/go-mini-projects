package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"net/http"
)

type Tunnel struct {
	w    io.Writer
	done chan struct{}
}

func makeSSHHandler(tunnel chan Tunnel) ssh.Handler {
	return func(session ssh.Session) {
		io.WriteString(session, "waiting for client...\n")

		t := <-tunnel

		io.WriteString(session, "client connected. copying...\n")

		n, err := io.Copy(t.w, session)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%d bytes written\n", n)
		close(t.done)

		io.WriteString(session, "transfer complete\n")
	}
}

func makeApiHandler(tunnel chan Tunnel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		t := Tunnel{
			w:    w,
			done: done,
		}

		log.Println("request is creating tunnel. waiting...")

		tunnel <- t

		log.Println("transferring...")

		<-done

		log.Println("transfer completed")
	}
}
