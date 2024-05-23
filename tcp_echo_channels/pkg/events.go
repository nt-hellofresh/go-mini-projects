package pkg

type InboundEvents struct {
	OnNewConnection      chan Transport
	OnNewMessageReceived chan []byte
}

type OutboundEvents struct {
	OnMessageSent chan []byte
}

type PeerEvents struct {
	In  InboundEvents
	Out OutboundEvents
}

func NewPeerEvents() PeerEvents {
	return PeerEvents{
		In: InboundEvents{
			OnNewConnection:      make(chan Transport),
			OnNewMessageReceived: make(chan []byte),
		},
		Out: OutboundEvents{
			OnMessageSent: make(chan []byte),
		},
	}
}