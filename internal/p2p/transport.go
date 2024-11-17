package p2p

import "net"

type Peer interface {
	net.Conn
	Send([]byte) error
	String() string
}

type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
}
