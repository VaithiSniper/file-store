package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

// MAX_DECODER_BUFFER_SIZE denotes the max size of the decoder buffer, and its value must be in sync with util.MaxAllowedDataPayloadSize
const MAX_DECODER_BUFFER_SIZE = 1024

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
	buf := make([]byte, MAX_DECODER_BUFFER_SIZE*2)
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
		var controlPayload ControlPayload
		decoder := gob.NewDecoder(bytes.NewReader(buf[:n]))
		if err := decoder.Decode(&controlPayload); err != nil {
			return fmt.Errorf("failed to decode control payload: %w", err)
		}
		m.Payload = controlPayload
	default:
		return fmt.Errorf("unsupported payload type: %T", m.Payload)
	}

	return nil
}
