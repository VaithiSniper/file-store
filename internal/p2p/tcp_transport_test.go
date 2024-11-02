package p2p

import (
	"file-store/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: NOHANDSHAKE,
		Decoder:       DefaultDecoder{},
	}
	tTransport := NewTCPTransport(tcpOpts, util.MessageChanBufferSize)

	assert.Equal(t, tTransport.TCPTransportOpts.ListenAddress, tcpOpts.ListenAddress)

	assert.Nil(t, tTransport.ListenAndAccept())

}
