package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

// MAX_DECODER_BUFFER_SIZE denotes the max size of the decoder buffer, and its value must be in sync with util.MaxAllowedDataPayloadSize
const MAX_DECODER_BUFFER_SIZE = 1024

type Codec interface {
	Encode(io.Writer, *Message) error
	Decode(io.Reader, *Message) error
}

type DefaultCodec struct{}

func (c *DefaultCodec) Encode(w io.Writer, msg *Message) error {
	// Helper lambda to write bytes into w
	var writeIntoW = func(b []byte) (int, error) {
		if size, err := w.Write(b); err != nil {
			return 0, fmt.Errorf("failed to write: %w", err)
		} else {
			return size, nil
		}
	}

	// Encode and write message type
	msgTypeBuf := []byte{byte(msg.Type)}
	if _, err := writeIntoW(msgTypeBuf); err != nil {
		return err
	}

	// Buffer to hold payload bytes
	var payloadBuf bytes.Buffer

	// Encode payload based on MessageType
	switch msg.Type {
	case DataMessageType:
		if payload, ok := msg.Payload.(DataPayload); ok {
			if err := gob.NewEncoder(&payloadBuf).Encode(payload); err != nil {
				return fmt.Errorf("failed to encode DataPayload Key: %w", err)
			}
		} else {
			return fmt.Errorf("invalid payload type for DataMessageType")
		}

	case ControlMessageType:
		if payload, ok := msg.Payload.(ControlPayload); ok {
			if err := gob.NewEncoder(&payloadBuf).Encode(payload); err != nil {
				return fmt.Errorf("failed to encode ControlPayload: %w", err)
			}
		} else {
			return fmt.Errorf("invalid payload type for ControlMessageType")
		}

	default:
		return fmt.Errorf("unsupported message type: %v", msg.Type)
	}

	// Write payload bytes to writer
	if _, err := writeIntoW(payloadBuf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (c *DefaultCodec) Decode(r io.Reader, msg *Message) error {
	// Helper lambda to read bytes from r
	var readFromR = func(buf []byte) (int, error) {
		if size, err := r.Read(buf); err != nil {
			return 0, fmt.Errorf("failed to read: %w", err)
		} else {
			return size, nil
		}
	}

	// Decode message type
	typeBuf := make([]byte, 1)
	if _, err := readFromR(typeBuf); err != nil {
		return err
	}
	msg.Type = MessageType(typeBuf[0])

	// Buffer for payload
	buf := make([]byte, MAX_DECODER_BUFFER_SIZE)
	n, err := readFromR(buf)
	if err != nil {
		return err
	}

	// Decode payload based on MessageType
	switch msg.Type {
	case DataMessageType:
		// Decode DataPayload
		var payload DataPayload
		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&payload); err != nil {
			return fmt.Errorf("failed to decode DataPayload key: %w", err)
		}
		msg.Payload = payload
		log.Printf("Received and decoded DataPayload -> %+v", msg.Payload)

	case ControlMessageType:
		// Decode ControlPayload
		var controlPayload ControlPayload
		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&controlPayload); err != nil {
			return fmt.Errorf("failed to decode ControlPayload: %w", err)
		}
		msg.Payload = controlPayload
		log.Printf("Received and decoded ControlPayload -> %+v", msg.Payload)

	default:
		return fmt.Errorf("unsupported payload type: %T", msg.Payload)
	}

	return nil
}
