package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type ControlMessage int

func (m ControlMessage) String() string {
	return [...]string{"STORE", "FETCH", "LIST", "EXIT"}[m]
}

type DataPayload struct {
	Key  string
	Data []byte
}

type Message struct {
	From    net.Addr
	Payload DataPayload
}

// ParseMessage decodes a message received from the network
func ParseMessage(msg Message) *Message {
	// Create a buffer from the received payload data
	buf := bytes.NewBuffer(msg.Payload.Data)

	// Create a new message to store the decoded data
	var decodedMsg Message
	if err := gob.NewDecoder(buf).Decode(&decodedMsg); err != nil {
		log.Printf("error decoding message: %+v", err)
	}

	// Return the decoded message with the original sender address
	decodedMsg.From = msg.From
	return &decodedMsg
}

// parseControlMessageType converts a string to a corresponding ControlMessage type.
func parseControlMessageType(str string) ControlMessage {
	switch str {
	case MESSAGE_FETCH_CONTROL_COMMAND.String():
		return MESSAGE_FETCH_CONTROL_COMMAND
	case MESSAGE_STORE_CONTROL_COMMAND.String():
		return MESSAGE_STORE_CONTROL_COMMAND
	case MESSAGE_LIST_CONTROL_COMMAND.String():
		return MESSAGE_LIST_CONTROL_COMMAND
	case MESSAGE_EXIT_CONTROL_COMMAND.String():
		return MESSAGE_EXIT_CONTROL_COMMAND
	default:
		return MESSAGE_UNKNOWN_CONTROL_COMMAND
	}
}

const (
	MESSAGE_STORE_CONTROL_COMMAND ControlMessage = iota
	MESSAGE_FETCH_CONTROL_COMMAND
	MESSAGE_LIST_CONTROL_COMMAND
	MESSAGE_EXIT_CONTROL_COMMAND
	MESSAGE_UNKNOWN_CONTROL_COMMAND
)
