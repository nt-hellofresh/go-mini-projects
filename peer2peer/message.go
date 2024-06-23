package main

import "net"

type RPC struct {
	From    net.Addr
	Payload []byte
}
