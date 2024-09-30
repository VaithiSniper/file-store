package p2p

import "net"

type ControlMessage int

const (
	MESSAGE_STORE_CONTROL_COMMAND ControlMessage = iota
	MESSAGE_FETCH_CONTROL_COMMAND
	MESSAGE_LIST_CONTROL_COMMAND
	MESSAGE_UNKNOWN_CONTROL_COMMAND
)

func (m ControlMessage) String() string {
	return [...]string{"STORE", "FETCH", "LIST"}[m]
}

type Message struct {
	From    net.Addr
	Payload []byte
}
