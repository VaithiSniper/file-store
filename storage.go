package main

import (
	"crypto/sha1"
	"encoding/hex"
	"file-store/file"
	"file-store/p2p"
	"fmt"
	"io"
	"strings"
)

const (
	DefaultChunkSize uint8 = 10
)

type PathTransformFunc func(string) string

var DefaultTransformFunc PathTransformFunc = func(key string) string {
	return key
}

var ContentAddressableTransformFunc PathTransformFunc = func(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	hashPath := strings.Join(chunkString(hashStr, DefaultChunkSize), "/")
	return hashPath
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts StoreOpts
}

var globalStore *Store

func getStoreInstance() *Store {
	if globalStore == nil {
		globalStore = createStoreWithDefaultOptions()
	}
	return globalStore
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s Store) writeFileToStorage(key string, r io.Reader) error {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
		FileMode: file.Default,
	}
	if err := f.WriteStream(r); err != nil {
		fmt.Println("Store Error: Error occurred while writing file to storage", err)
		return err
	}
	return nil
}

func (s Store) readFileFromStorage(key string) ([]byte, error) {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.ReadFile()
}

func (s Store) deleteFileFromStorage(key string) error {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.DeleteFile()
}

func (s Store) existsInStorage(key string) bool {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.Exists()
}

func (s Store) parseControlMessageType(str string) p2p.ControlMessage {
	switch str {
	case p2p.MESSAGE_FETCH_CONTROL_COMMAND.String():
		return p2p.MESSAGE_FETCH_CONTROL_COMMAND
	case p2p.MESSAGE_STORE_CONTROL_COMMAND.String():
		return p2p.MESSAGE_STORE_CONTROL_COMMAND
	case p2p.MESSAGE_LIST_CONTROL_COMMAND.String():
		return p2p.MESSAGE_LIST_CONTROL_COMMAND
	default:
		return p2p.MESSAGE_UNKNOWN_CONTROL_COMMAND
	}
}

// chunkString chunks the given hex string s into blocks of fixed size uint8 blockSize
func chunkString(s string, blockSize uint8) []string {
	var chunks []string
	i := uint8(0)
	strLen := uint8(len(s))
	for i = 0; i < strLen; i += blockSize {
		end := i + blockSize
		if end > strLen {
			end = strLen
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func createStoreWithDefaultOptions() *Store {
	opts := StoreOpts{PathTransformFunc: ContentAddressableTransformFunc}
	store := NewStore(opts)
	return store
}
