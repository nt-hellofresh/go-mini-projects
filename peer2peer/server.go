package main

import (
	"encoding/gob"
	"log"
	"sync"
)

type ServerOpts struct {
	Transport        Transport
	ConnectAddresses []string
}

type Server struct {
	ServerOpts

	peerLock sync.Mutex
	peers    map[string]Peer
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		peers:      make(map[string]Peer),
	}
}

func (s *Server) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if err := s.connectToPeers(); err != nil {
		return err
	}

	return s.loop()
}

func (s *Server) connectToPeers() error {
	for _, addr := range s.ConnectAddresses {
		if err := s.Transport.Dial(addr); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) broadcast(data []byte) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	rpc := RPC{
		Payload: data,
	}

	for _, p := range s.peers {
		if err := gob.NewEncoder(p).Encode(rpc); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) onConnect(peer Peer) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.Addr()] = peer
}

func (s *Server) loop() error {
	for {
		select {
		case peer := <-s.Transport.OnConnect():
			s.onConnect(peer)
		case rpc := <-s.Transport.OnMessage():
			log.Printf("%v received message from %v\n", s.Transport.Addr(), rpc.From)
		}
	}
}
