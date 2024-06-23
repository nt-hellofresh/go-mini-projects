package main

import "net"

type Peer interface {
	Addr() string
	Close() error
	Write([]byte) (int, error)
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func (p *TCPPeer) Addr() string {
	return p.conn.RemoteAddr().String()
}

func (p *TCPPeer) Write(data []byte) (int, error) {
	return p.conn.Write(data)
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}
