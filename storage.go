package main

import (
	"crypto/sha1"
	"encoding/hex"
	"file-store/file"
	"file-store/p2p"
	"file-store/util"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
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
	hashPath := strings.Join(util.ChunkString(hashStr, DefaultChunkSize), "/")
	return hashPath
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
	MessageFormat     p2p.MessageFormat
}

type Store struct {
	StoreOpts StoreOpts
}

var globalStore *Store

// createStoreWithDefaultOptions initializes a Store with default options using a content-addressable path transform function.
func createStoreWithDefaultOptions() *Store {
	opts := StoreOpts{PathTransformFunc: ContentAddressableTransformFunc, MessageFormat: p2p.JSONFormat{}}
	store := Store{
		StoreOpts: opts,
	}
	return &store
}

// getStoreInstance returns a singleton instance of Store. If the instance doesn't exist, it creates one with default options.
func getStoreInstance() *Store {
	if globalStore == nil {
		globalStore = createStoreWithDefaultOptions()
	}
	return globalStore
}

// --------------------------------------------------------------  CONTROL PLANE --------------------------------------------------------------
func (s Store) setupHyperStoreServer() {
	// To wait on goroutine
	var wg sync.WaitGroup

	// Prepare server with opts
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: p2p.NOHANDSHAKE,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        onPeerSuccess,
	}
	tTransport := p2p.NewTCPTransport(tcpOpts)

	// Start listening for incoming connections
	fmt.Println("Starting to listen and accept connections...")
	if err := tTransport.ListenAndAccept(); err != nil {
		log.Fatalln("Error listening and accepting connections:", err)
	}

	wg.Add(1)
	// Start read loop
	go handlePeerRead(tTransport, &wg)
	// Wait until peer exists
	wg.Wait()
}

func handlePeerRead(tTransport *p2p.TCPTransport, wg *sync.WaitGroup) {
	defer wg.Done()
	var msgCount uint32 = 0
	for {
		msg := <-tTransport.Consume()
		controlMessage, err := msg.ParseMessage()
		if err != nil || controlMessage == p2p.MESSAGE_UNKNOWN_CONTROL_COMMAND {
			log.Fatalf("Error parsing message: %+v", err)
		}
		if controlMessage == p2p.MESSAGE_EXIT_CONTROL_COMMAND {
			fmt.Printf("Finished reading %d messages from peer %s\n", msgCount, msg.From)
			break
		}
		msgCount++
	}
}

// --------------------------------------------------------------  CONTROL PLANE --------------------------------------------------------------

// --------------------------------------------------------------  FILE HANDLING --------------------------------------------------------------
// handleFileWrite writes the content from the given io.Reader to a file specified by the key within the storage system.
func (s Store) handleFileWrite(key string, r io.Reader) error {
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

// handleFileRead reads the file identified by the given key and returns its content as a byte slice.
func (s Store) handleFileRead(key string) ([]byte, error) {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.ReadFile()
}

// handleFileDelete deletes the file identified by the given key within the storage system.
func (s Store) handleFileDelete(key string) error {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.DeleteFile()
}

// existsInStorage checks if a file identified by the given key exists in the storage system.
func (s Store) existsInStorage(key string) bool {
	pathname := s.StoreOpts.PathTransformFunc(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.Exists()
}

// --------------------------------------------------------------  END OF FILE HANDLING --------------------------------------------------------------
