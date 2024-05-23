package main

import (
	"log"
	"tcp_echo_channels/pkg"
	"time"
)

func main() {
	events := pkg.NewPeerEvents()
	connector := pkg.NewTransportConnector(":3000", events.In.OnNewConnection)

	opts := pkg.PeerOpts{
		Events:         events,
		Connector:      connector,
		EchoServerMode: true,
	}
	peer := pkg.NewPeer(opts)

	go peer.Start()

	time.Sleep(1 * time.Second)

	RunClient()
}

func RunClient() {
	events := pkg.NewPeerEvents()
	connector := pkg.NewTransportConnector(":4001", events.In.OnNewConnection)

	opts := pkg.PeerOpts{
		Events:         events,
		Connector:      connector,
		EchoServerMode: false,
	}
	client := pkg.NewPeer(opts)

	go client.Start()

	err := connector.Dial(":3000")

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	client.Broadcast([]byte("Foo"))

	time.Sleep(1 * time.Second)
}