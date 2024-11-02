package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type ControlMessage int

const (
	MESSAGE_STORE_CONTROL_COMMAND ControlMessage = iota
	MESSAGE_FETCH_CONTROL_COMMAND
	MESSAGE_LIST_CONTROL_COMMAND
	MESSAGE_EXIT_CONTROL_COMMAND
	MESSAGE_UNKNOWN_CONTROL_COMMAND
)

func (m ControlMessage) String() string {
	return [...]string{"STORE", "FETCH", "LIST", "EXIT"}[m]
}

// MessageType denotes the type of message received from an enum list
type MessageType byte

const (
	DataMessageType MessageType = iota
	ControlMessageType
)

// DataPayload represents file transfer data
type DataPayload struct {
	Key  string
	Data []byte
}

// ControlPayload represents control messages
type ControlPayload struct {
	Command ControlMessage
	//Args    map[string]string
}

type Message struct {
	Type    MessageType
	From    net.Addr
	Payload interface{}
}

// ParseMessage decodes a message received from the network
func ParseMessage(msg Message) *Message {
	// Create a new message to store the decoded data
	var decodedMsg Message
	// Copy over the original sender address
	decodedMsg.From = msg.From
	decodedMsg.Type = msg.Type

	switch msg.Type {
	case DataMessageType:
		ParseDataMessage(msg, &decodedMsg)
	case ControlMessageType:
		// TODO: Implement ControlPayload handling
		ParseControlMessage(msg, &decodedMsg)
	default: // Do nothing
	}

	return &decodedMsg
}

// ParseControlMessage handles parsing of control payloads into decodedMsg
func ParseControlMessage(msg Message, decodedMsg *Message) *Message {
	decodedMsg.Payload = ControlPayload{Command: msg.Payload.(ControlPayload).Command}
	return decodedMsg
}

// ParseDataMessage handles parsing of data payloads into decodedMsg
func ParseDataMessage(msg Message, decodedMsg *Message) *Message {
	buf := bytes.NewBuffer(msg.Payload.(DataPayload).Data)
	if err := gob.NewDecoder(buf).Decode(&decodedMsg); err != nil {
		log.Printf("error decoding data message: %+v", err)
		return nil
	}
	return decodedMsg
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
