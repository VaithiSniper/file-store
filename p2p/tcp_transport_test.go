package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: NOHANDSHAKE,
		Decoder:       DefaultDecoder{},
	}
	tTransport := NewTCPTransport(tcpOpts)

	assert.Equal(t, tTransport.TCPTransportOpts.ListenAddress, tcpOpts.ListenAddress)

	assert.Nil(t, tTransport.ListenAndAccept())

}
