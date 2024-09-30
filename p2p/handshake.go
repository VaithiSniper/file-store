package p2p

import "errors"

var ErrInvalidHandshake = errors.New("Invalid handshake, couldn't verify peer!")

type doHandshake func(Peer) error

func NOHANDSHAKE(Peer) error {
	return nil
}
