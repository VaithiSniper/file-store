package p2p

import (
	"encoding/gob"
	"fmt"
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
	switch m.Type {
	case DataMessageType:
		payload := DataPayload{
			Data: buf[:n],
		}
		m.Payload = payload
	case ControlMessageType:
		// TODO: Handle ControlMessageType decoding
	default:
		return fmt.Errorf("unsupported payload type: %T", m.Payload)
	}

	return nil
}
