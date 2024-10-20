package p2p

import (
	"file-store/internal/util"
	"fmt"
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

type Message struct {
	From    net.Addr
	Payload []byte
}

func (msg *Message) ParseMessage() (ControlMessage, error) {
	var controlMessage ControlMessage
	str, err := util.SafeByteToString(msg.Payload)
	if err == nil {
		controlMessage = parseControlMessageType(str)
		if controlMessage == MESSAGE_UNKNOWN_CONTROL_COMMAND {
			err = fmt.Errorf("unknown control message: %s", str)
		}
	}
	return controlMessage, err
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
