package util

import (
	"encoding/gob"
	"file-store/internal/p2p"
	"net"
)

func RegisterGobTypes() {
	// Register types with gob
	gob.Register(&net.TCPAddr{})
	gob.Register(p2p.Message{})
	gob.Register(p2p.DataPayload{})
	gob.Register(p2p.ControlPayload{})
	gob.Register(map[string]string{})
}

// STORE_ACTION is for values that denotes actions that can be performed on the Store by the user
type STORE_ACTION int

// MAP_ACTION is for the possible operations that be performed on a map
type MAP_ACTION int
