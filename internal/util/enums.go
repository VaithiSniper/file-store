package util

// List of actions that can be performed by user on the Store
const (
	STORE_ACTION_STORE STORE_ACTION = iota
	STORE_ACTION_DELETE
)

// Map opertaions that are possible
const (
	MAP_GET_ELEMENT MAP_ACTION = iota
	MAP_UPSERT_ELEMENT
	MAP_DELETE_ELEMENT
)
