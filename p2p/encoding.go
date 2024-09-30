package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct {
}

type DefaultDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, m *Message) error {
	return gob.NewDecoder(r).Decode(m)
}

func (dec DefaultDecoder) Decode(r io.Reader, m *Message) error {
	buf := make([]byte, 1024*2)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	m.Payload = buf[:n]
	return nil
}
