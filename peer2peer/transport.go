package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Transport interface {
	Addr() string
	ListenAndAccept() error
	OnMessage() <-chan RPC
	OnConnect() <-chan Peer
	Dial(address string) error
}

func defaultHandshake(conn net.Conn) error {
	return nil
}

type TCPTransportOpts struct {
	listenAddress string
	handshake     func(net.Conn) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener    net.Listener
	rpcCh       chan RPC
	connectedCh chan Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC, 1024),
		connectedCh:      make(chan Peer, 1024),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Println("Listening on", t.listenAddress)

	return nil
}

func (t *TCPTransport) Addr() string {
	return t.listenAddress
}

func (t *TCPTransport) OnConnect() <-chan Peer {
	return t.connectedCh
}

func (t *TCPTransport) OnMessage() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)

	defer func() {
		log.Println("Closing connection to", conn.RemoteAddr())
		conn.Close()
	}()

	if err := t.handshake(conn); err != nil {
		return
	}

	t.connectedCh <- peer

	//Read loop
	for {
		msg := &RPC{}

		err := gob.NewDecoder(conn).Decode(msg)

		if err != nil {
			fmt.Printf("Failed to decode message: %v\n", err)
			return
		}

		msg.From = conn.RemoteAddr()
		t.rpcCh <- *msg
	}
}
