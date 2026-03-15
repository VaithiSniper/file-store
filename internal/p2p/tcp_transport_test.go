package p2p

import (
	"file-store/internal/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: NOHANDSHAKE,
		Codec:         &DefaultCodec{},
	}
	tTransport := NewTCPTransport(tcpOpts, constants.MessageChanBufferSize)

	assert.Equal(
		t, tTransport.TCPTransportOpts.ListenAddress, tcpOpts.ListenAddress,
	)

	assert.Nil(t, tTransport.ListenAndAccept())

}
