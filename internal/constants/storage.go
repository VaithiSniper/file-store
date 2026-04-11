package constants

import "time"

const (
	DefaultNodeName               string = "alpha-node-1"
	DefaultBaseStorageLocation    string = "storage"
	DefaultListenAddress          string = ":5000"
	DefaultAPIServerListenAddress string = ":8080"
)

const (
	FetchMessageResponseTimeout = 15 * time.Second
)
