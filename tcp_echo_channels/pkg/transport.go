package pkg

import (
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type Transport struct {
	conn io.ReadWriteCloser
	PID  uuid.UUID
}

func NewTransport(conn io.ReadWriteCloser) Transport {
	return Transport{
		conn: conn,
		PID:  uuid.New(),
	}
}

func (tr Transport) ReadLoop(rxChan chan []byte) {
	for {
		buf := make([]byte, 1024)

		n, err := tr.conn.Read(buf)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%v Receiving\n", tr.PID)
		rxChan <- buf[:n]
	}
}

func (tr Transport) WriteLoop(txChan chan []byte) {
	defer tr.conn.Close()
	for {
		select {
		case msg := <-txChan:
			log.Printf("%v Transmitting\n", tr.PID)
			err := tr.send(msg)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (tr Transport) send(data []byte) error {
	_, err := tr.conn.Write(data)

	if err != nil {
		return err
	}

	return nil

}

type TransportConnector interface {
	Dial(address string) error
	AcceptConnectionsLoop()
}

type TCPTransportConnector struct {
	listenAddress   string
	onNewConnection chan Transport
}

func NewTransportConnector(address string, newConnCh chan Transport) *TCPTransportConnector {
	return &TCPTransportConnector{
		listenAddress:   address,
		onNewConnection: newConnCh,
	}
}

func (ttc *TCPTransportConnector) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return err
	}

	ttc.onNewConnection <- NewTransport(conn)

	return nil
}

func (ttc *TCPTransportConnector) AcceptConnectionsLoop() {
	ln, err := net.Listen("tcp", ttc.listenAddress)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		ttc.onNewConnection <- NewTransport(conn)
	}
}