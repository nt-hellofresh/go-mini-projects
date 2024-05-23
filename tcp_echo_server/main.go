package main

import (
	"tcp_echo_server/pkg"
	"time"
)

func main() {
	server := &pkg.EchoServer{}
	go server.Start()

	time.Sleep(1 * time.Second)

	pkg.RunClient()
}
