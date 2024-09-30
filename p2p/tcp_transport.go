package p2p

import (
	"fmt"
	"net"
	"reflect"
)

type TCPTransport struct {
	TCPTransportOpts
	listener    net.Listener
	messageChan chan Message
}

type TCPTransportOpts struct {
	ListenAddress string
	HandshakeFunc doHandshake
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPPeer struct {
	conn net.Conn
	// for tcp-dial => true, for tcp-accept => false
	isOutbound bool
}

func NewTCPPeer(conn net.Conn, isOutbound bool) *TCPPeer {
	return &TCPPeer{
		conn:       conn,
		isOutbound: isOutbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		messageChan:      make(chan Message),
	}
}

// Consume implements of the Transport interface, and returns read-only channel ref
func (t *TCPTransport) Consume() <-chan Message {
	return t.messageChan
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.accept()
	return nil
}

func (t *TCPTransport) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			err := fmt.Errorf("TCP Error: Error while accepting connection: %s\n", err)
			fmt.Println(err.Error())
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Println("Dropping peer connection, connection ending...")
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	fmt.Println("TCP Error: New connection from peer: " + peer.conn.RemoteAddr().String())

	// Perform handshake and authenticate peer
	if err = t.HandshakeFunc(peer); err != nil {
		_ = peer.Close()
		fmt.Println("TCP Error: Error while handshaking, closing connection to " + peer.conn.RemoteAddr().String())
		return
	}

	// Call the onPeer on the peer ifc
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Println("Error while using onPeer, terminating connection")
			_ = peer.Close()
			return
		}
	}

	fmt.Println("Entering read loop..." + peer.conn.RemoteAddr().String())
	// Once authenticated, read messages in read loop
	msg := Message{}
	for {
		// Decode message from conn to msg
		err := t.Decoder.Decode(conn, &msg)
		//  TODO: Handle abrupt peer disconnect during onPeer func, since it comes to read loop at that point
		fmt.Printf("Type of error +%v\n", reflect.TypeOf(err))
		if err != nil {
			fmt.Println("TCP Error: Error decoding bytes from connection")
			continue
		}
		// Else, message exists and is parsed
		msg.From = conn.RemoteAddr() // Picks up socket wrt current connection, but ideally need IP of other peer
		t.messageChan <- msg
	}
}
