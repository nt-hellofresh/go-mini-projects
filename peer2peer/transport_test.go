package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		listenAddress: ":3333",
		handshake:     defaultHandshake,
	}
	tt := NewTCPTransport(opts)

	assert.Equal(t, opts.listenAddress, tt.listenAddress)
	assert.Nil(t, tt.ListenAndAccept())
}
