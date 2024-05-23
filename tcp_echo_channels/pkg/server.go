package pkg

import (
	"fmt"
	"log"
)

type Peer struct {
	events     PeerEvents
	connector  TransportConnector
	echoServer bool
}

type PeerOpts struct {
	Events         PeerEvents
	Connector      TransportConnector
	EchoServerMode bool
}

func NewPeer(opts PeerOpts) *Peer {
	return &Peer{
		events:     opts.Events,
		connector:  opts.Connector,
		echoServer: opts.EchoServerMode,
	}
}

func (p *Peer) Start() {
	go p.subscribeToEvents()

	p.connector.AcceptConnectionsLoop()
}

func (p *Peer) subscribeToEvents() {
	for {
		select {
		case conn := <-p.events.In.OnNewConnection:
			log.Printf("New connection %v", conn.PID)
			go conn.ReadLoop(p.events.In.OnNewMessageReceived)
			go conn.WriteLoop(p.events.Out.OnMessageSent)
		case buf := <-p.events.In.OnNewMessageReceived:
			log.Printf("Peer received '%v'", string(buf))

			if p.echoServer {
				data := string(buf)
				respData := []byte(fmt.Sprintf("You sent me: %v", data))
				p.Broadcast(respData)
			}
		}
	}
}

func (p *Peer) Broadcast(message []byte) {
	// can end up blocking if no clients connected
	p.events.Out.OnMessageSent <- message
}
