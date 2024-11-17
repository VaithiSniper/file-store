package p2p

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type ControlMessage int

const (
	MESSAGE_STORE_CONTROL_COMMAND ControlMessage = iota
	MESSAGE_FETCH_CONTROL_COMMAND
	MESSAGE_FETCH_RESPONSE_CONTROL_COMMAND
	MESSAGE_LIST_CONTROL_COMMAND
	MESSAGE_EXIT_CONTROL_COMMAND
	MESSAGE_UNKNOWN_CONTROL_COMMAND
)

func (m ControlMessage) String() string {
	return [...]string{"STORE", "FETCH", "FETCH_RESPONSE", "LIST", "EXIT"}[m]
}

// MessageType denotes the type of message received from an enum list
type MessageType byte

const (
	DataMessageType MessageType = iota
	ControlMessageType
)

func (m MessageType) String() string {
	return [...]string{"DATA", "CONTROL"}[m]
}

// DataPayload represents file transfer data
type DataPayload struct {
	Key  string
	Data []byte
}

func (d *DataPayload) String() string {
	return fmt.Sprintf("DataPayload containing Key=%s and Data=%+v", d.Key, d.Data)
}

// ControlPayload represents control messages
type ControlPayload struct {
	Command ControlMessage
	Args    map[string]string
}

func (c *ControlPayload) String() string {
	return fmt.Sprintf("ControlPayload containing Command=%s and Args=%+v", c.Command, c.Args)
}

type Message struct {
	Type    MessageType
	From    net.Addr
	Payload interface{}
}

func (m *Message) String() string {
	return fmt.Sprintf("Message containing Type=%s, From=%s and Payload=%+v", m.Type, m.From, m.Payload)
}

// ConstructFetchResponseMessage constructs and return MESSAGE_FETCH_RESPONSE_CONTROL_COMMAND message based on whether the file was found or not
func ConstructFetchResponseMessage(fileExists bool) Message {
	return Message{
		Type: ControlMessageType,
		Payload: &ControlPayload{
			Command: MESSAGE_FETCH_RESPONSE_CONTROL_COMMAND,
			Args: map[string]string{
				"file_exists": strconv.FormatBool(fileExists),
			},
		},
	}
}

// ParseMessage decodes a message received from the network
func ParseMessage(msg Message) *Message {
	var decodedMsg Message
	decodedMsg.From = msg.From
	decodedMsg.Type = msg.Type

	switch msg.Type {
	case DataMessageType:
		if decodedMsgPayload, ok := msg.Payload.(DataPayload); ok {
			decodedMsg.Payload = decodedMsgPayload
		} else {
			log.Printf("Error parsing message into DataMessageType")
			return nil
		}
	case ControlMessageType:
		if decodedMsgPayload, ok := msg.Payload.(ControlPayload); ok {
			decodedMsg.Payload = decodedMsgPayload
		} else {
			log.Printf("Error parsing message into ControlMessageType")
			return nil
		}
	default:
		log.Printf("Unknown message type: %d", msg.Type)
		return nil
	}

	return &decodedMsg
}
