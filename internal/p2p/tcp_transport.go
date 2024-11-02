package p2p

import (
	"fmt"
	"io"
	"net"
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
	net.Conn
	// for tcp-dial => true, for tcp-accept => false
	isOutbound bool
}

func NewTCPPeer(conn net.Conn, isOutbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:       conn,
		isOutbound: isOutbound,
	}
}

// Send implements the Peer interface, and send the given msg bytes to that peer
func (p *TCPPeer) Send(msg []byte) error {
	_, err := p.Conn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		messageChan:      make(chan Message),
	}
}

// Consume implements the Transport interface, and returns read-only channel ref
func (t *TCPTransport) Consume() <-chan Message {
	return t.messageChan
}

// Close implements the Transport interface, closes the transport channel and returns err
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// ListenAndAccept implements the Transport interface, listens on t.ListenAddress for incoming connections
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.accept()
	return nil
}

// Dial implements the Transport interface, dials out to bootstrap nodes and sets them up as part of the network
func (t *TCPTransport) Dial(nodeAddr string) error {
	conn, err := net.Dial("tcp", nodeAddr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			err := fmt.Errorf("TCP Error: Error while accepting connection: %s\n", err)
			fmt.Println(err.Error())
		}
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, isOutbound bool) {
	var err error
	defer func() {
		fmt.Println("Dropping peer connection, connection ending...")
		err := conn.Close()
		if err != nil {
			fmt.Println(fmt.Errorf("Error while closing connection: %s\n", err))
			return
		}
	}()

	peer := NewTCPPeer(conn, isOutbound)
	fmt.Println("New connection from peer: " + peer.RemoteAddr().String())

	// Perform handshake and authenticate peer
	if err = t.HandshakeFunc(peer); err != nil {
		_ = peer.Close()
		fmt.Println("TCP Error: Error while handshaking, closing connection to " + peer.RemoteAddr().String())
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

	fmt.Println("Entering read loop..." + peer.RemoteAddr().String())
	// Once authenticated, read messages in read loop
	msg := Message{}
	for {
		// Decode message from conn to msg
		err := t.Decoder.Decode(conn, &msg)
		//  TODO: Handle abrupt peer disconnect during onPeer func, since it comes to read loop at that point
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Peer %s disconnected\n", peer.RemoteAddr().String())
			} else {
				fmt.Printf("Error decoding message from %s: %v\n", peer.RemoteAddr().String(), err)
			}
			return
		}

		// Set the sender address and forward the message
		msg.From = conn.RemoteAddr()
		select {
		case t.messageChan <- msg:
			// Message forwarded successfully
		default:
			fmt.Printf("Warning: Message channel full, dropping message from %s\n",
				peer.RemoteAddr().String())
		}
	}
}
