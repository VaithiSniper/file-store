package main

import (
	"crypto/sha1"
	"encoding/hex"
	"file-store/internal/file"
	"file-store/internal/p2p"
	"file-store/internal/util"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	DefaultChunkSize uint8 = 10
)

// PathTransformFunc is the type of any function that takes in a key and base storage location and returns the complete path to store the file
type PathTransformFunc func(baseStorageLocation string, key string) string

// DefaultTransformFunc is an implementation of PathTransformFunc that just preserves the original key
var DefaultTransformFunc PathTransformFunc = func(baseStorageLocation string, key string) string {
	return key
}

// ContentAddressableTransformFunc is an implementation of PathTransformFunc that uses sha1 hashing to generate the path
var ContentAddressableTransformFunc PathTransformFunc = func(baseStorageLocation string, key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	hashPath := strings.Join(util.ChunkString(hashStr, DefaultChunkSize), "/")
	return hashPath
}

type StoreOpts struct {
	PathTransformFunc   PathTransformFunc
	MessageFormat       p2p.MessageFormat
	BaseStorageLocation string
}

type Store struct {
	StoreOpts StoreOpts
	Transport p2p.Transport
}

var globalStore *Store

// createStoreWithDefaultOptions initializes a Store with default options using a content-addressable path transform function.
func createStoreWithDefaultOptions() *Store {
	// Prepare Transport with opts
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: p2p.NOHANDSHAKE,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        onPeerSuccess,
	}
	tTransport := p2p.NewTCPTransport(tcpOpts)
	// Prepare Store with opts
	opts := StoreOpts{PathTransformFunc: ContentAddressableTransformFunc, MessageFormat: p2p.JSONFormat{}, BaseStorageLocation: util.DefaultBaseStorageLocation}
	store := Store{
		StoreOpts: opts,
		Transport: tTransport,
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

func (s *Store) setupHyperStoreServer() {
	var wg sync.WaitGroup

	// Start listening for incoming connections
	log.Println("Starting to listen and accept connections...")
	if err := s.Transport.ListenAndAccept(); err != nil {
		log.Fatalln("Error listening and accepting connections:", err)
	}
	log.Printf("Listening on %v", 5000)

	wg.Add(1)
	// Start read loop
	go s.handlePeerRead(&wg)
	// Wait until peer exists
	wg.Wait()
}

func (s *Store) teardownHyperStoreServer() {
	log.Println("Hyperstore stopped due to user STOP action")
	// Terminate all connections

	// Remove the base path entirely
	err := os.RemoveAll(s.StoreOpts.BaseStorageLocation)
	if err != nil {
		log.Fatalln("Error while tearing down:", err)
	}
}

func (s *Store) handlePeerRead(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		log.Println("Shutting down peer read due to peer QUIT")
		time.Sleep(time.Second * 3)

		s.Transport.Close()
	}()

	var msgCount uint32 = 0
	for {
		msg := <-s.Transport.Consume()
		controlMessage, err := msg.ParseMessage()
		if controlMessage == p2p.MESSAGE_EXIT_CONTROL_COMMAND {
			fmt.Printf("Finished reading %d messages from peer %s\n", msgCount, msg.From)
			s.teardownHyperStoreServer()
			break
		}
		switch controlMessage {
		case p2p.MESSAGE_STORE_CONTROL_COMMAND:
			fmt.Printf("Received STORE_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		case p2p.MESSAGE_LIST_CONTROL_COMMAND:
			fmt.Printf("Received LIST_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		case p2p.MESSAGE_FETCH_CONTROL_COMMAND:
			fmt.Printf("Received FETCH_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		default:
			log.Fatalf("Error parsing message: %+v from ---- %s\n", err, msg.From.String())
		}

		msgCount++
	}
}

// --------------------------------------------------------------  END OF CONTROL PLANE --------------------------------------------------------------

// --------------------------------------------------------------  FILE HANDLING --------------------------------------------------------------

// generatePath generates and returns a path to store a file with given key
func (s *Store) generatePath(key string) string {
	hashPath := s.StoreOpts.PathTransformFunc(key, s.StoreOpts.BaseStorageLocation)
	return path.Join(s.StoreOpts.BaseStorageLocation, hashPath)
}

// handleFileWrite writes the content from the given io.Reader to a file specified by the key within the storage system.
func (s *Store) handleFileWrite(key string, r io.Reader) error {
	pathname := s.generatePath(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
		FileMode: util.Default,
	}
	if err := f.WriteStream(r); err != nil {
		fmt.Println("Store Error: Error occurred while writing file to storage", err)
		return err
	}
	return nil
}

// handleFileRead reads the file identified by the given key and returns its content as a byte slice.
func (s *Store) handleFileRead(key string) ([]byte, error) {
	pathname := s.generatePath(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.ReadFile()
}

// handleFileDelete deletes the file identified by the given key within the storage system.
func (s *Store) handleFileDelete(key string) error {
	pathname := s.generatePath(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.DeleteFile()
}

// existsInStorage checks if a file identified by the given key exists in the storage system.
func (s *Store) existsInStorage(key string) bool {
	pathname := s.generatePath(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
	}
	return f.Exists()
}

// --------------------------------------------------------------  END OF FILE HANDLING --------------------------------------------------------------
