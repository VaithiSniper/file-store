package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
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

func (s *Store) OnPeer(p p2p.Peer) error {
	s.PeerLock.Lock()
	defer s.PeerLock.Unlock()

	log.Printf("Adding peer %s to PeerMap\n", p.RemoteAddr())
	s.PeerMap[p.RemoteAddr().String()] = p

	return nil
}

func onPeerFailure(peer p2p.Peer) error {
	return fmt.Errorf("error occuring")
}

func onPeerSuccess(peer p2p.Peer) error {
	return nil
}

func onPeerAbruptPeerCloseFailure(peer p2p.Peer) error {
	return peer.Close()
}

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
	hashPath := strings.Join(util.ChunkString(hashStr, util.DefaultChunkSize), "/")
	return hashPath
}

type StoreOpts struct {
	ListenAddress       string
	PathTransformFunc   PathTransformFunc
	MessageFormat       p2p.MessageFormat
	BaseStorageLocation string
	BootstrapNodes      []string
}

type Store struct {
	StoreOpts StoreOpts
	Transport p2p.Transport
	PeerLock  sync.Mutex
	PeerMap   map[string]p2p.Peer
}

var globalStore *Store

// createStoreWithDefaultOptions initializes a Store with default options using a content-addressable path transform function.
func createStoreWithDefaultOptions(listenAddress string, bootstrapNodes []string, fileStorageBasePath string) *Store {
	// Prepare Transport with opts
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddress,
		HandshakeFunc: p2p.NOHANDSHAKE,
		Decoder:       p2p.DefaultDecoder{},
	}
	tTransport := p2p.NewTCPTransport(tcpOpts, util.MessageChanBufferSize)
	// Prepare Store with opts
	opts := StoreOpts{
		ListenAddress:       listenAddress,
		PathTransformFunc:   ContentAddressableTransformFunc,
		MessageFormat:       p2p.JSONFormat{},
		BaseStorageLocation: fileStorageBasePath,
		BootstrapNodes:      bootstrapNodes,
	}
	store := Store{
		StoreOpts: opts,
		Transport: tTransport,
		PeerMap:   make(map[string]p2p.Peer),
	}
	// Set onPeer on Transport to use Store's onPeer method
	tTransport.OnPeer = store.OnPeer
	return &store
}

// --------------------------------------------------------------  CONTROL PLANE --------------------------------------------------------------

// getStoreInstance returns a singleton instance of Store. If the instance doesn't exist, it creates one with provided params.
func getStoreInstance(listenAddress string, bootstrapNodes []string, fileStorageBasePath string) *Store {
	if globalStore == nil {
		globalStore = createStoreWithDefaultOptions(listenAddress, bootstrapNodes, fileStorageBasePath)
	}
	return globalStore
}

// bootstrapNetwork with improved error handling and synchronization
func (s *Store) bootstrapNetwork() error {
	var wg sync.WaitGroup
	var errors []error
	var errorsMux sync.Mutex

	for _, nodeAddr := range s.StoreOpts.BootstrapNodes {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			if err := s.Transport.Dial(addr); err != nil {
				errorsMux.Lock()
				errors = append(errors, fmt.Errorf("failed to dial %s: %w", addr, err))
				errorsMux.Unlock()
			}
		}(nodeAddr)
	}

	wg.Wait()

	if len(errors) > 0 {
		// Return first error or combine them as needed
		return fmt.Errorf("bootstrap errors: %v", errors)
	}
	return nil
}

// setupHyperStoreServer starts the Store on provided ListenAddress
func (s *Store) setupHyperStoreServer() {
	var wg sync.WaitGroup

	// Start listening for incoming connections
	log.Println("Starting to listen and accept connections...")
	if err := s.Transport.ListenAndAccept(); err != nil {
		log.Fatalln("Error listening and accepting connections:", err)
	}
	log.Printf("Listening on %v", s.StoreOpts.ListenAddress)

	// Bootstrapping network with predefined nodes
	if len(s.StoreOpts.BootstrapNodes) == 0 {
		log.Println("No bootstrap nodes were specified.")
	} else {
		err := s.bootstrapNetwork()
		if err != nil {
			log.Println("Error while bootstrapping network:", err)
		}
	}

	wg.Add(1)
	// Start read loop
	go s.handlePeerRead(&wg)
	// Wait until peer exists
	wg.Wait()
}

// teardownHyperStoreServer terminates any existing connections, cleans up data/db and stops the Store
func (s *Store) teardownHyperStoreServer() {
	log.Println("Hyperstore stopped due to user STOP action")
	// Terminate all connections

	// Remove the base path entirely
	err := os.RemoveAll(s.StoreOpts.BaseStorageLocation)
	if err != nil {
		log.Fatalln("Error while tearing down:", err)
	}
}

// handlePeerRead causes the peer goes into a read loop where it reads from the msg channel
func (s *Store) handlePeerRead(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		log.Println("Shutting down peer read due to peer QUIT")
		time.Sleep(time.Second * 3)

		s.Transport.Close()
	}()

	var msgCount uint32 = 0
	var toRead = true
	// TODO: Set value of toRead to false somewhere in loop based on EXIT Control Message
	for toRead {
		msg := <-s.Transport.Consume()
		parsedMsg := p2p.ParseMessage(msg)

		var validateMessageParseDebug = func(parsedMsg *p2p.Message) {
			switch parsedMsg.Type {
			case p2p.DataMessageType:
				payload := parsedMsg.Payload.(p2p.DataPayload)
				log.Printf("Received data from %s: Key=%s, Data=%q",
					parsedMsg.From, payload.Key, string(payload.Data))
			case p2p.ControlMessageType:
				payload := parsedMsg.Payload.(p2p.ControlPayload)
				log.Printf("Received control from %s: Command=%s, CommandNumber=%d",
					parsedMsg.From, payload.Command, payload.Command)
				switch payload.Command {
				case p2p.MESSAGE_STORE_CONTROL_COMMAND:
					sender, senderExists := s.PeerMap[parsedMsg.From.String()]
					if senderExists {
						log.Printf("Sender %s exists in peerMap", sender.RemoteAddr())
					} else {
						log.Printf("Sender %s does not exist in peerMap", sender)
					}
				case p2p.MESSAGE_EXIT_CONTROL_COMMAND:
					sender, senderExists := s.PeerMap[parsedMsg.From.String()]
					if senderExists {
						log.Printf("Sender %s exists in peerMap", sender.RemoteAddr())
						toRead = false
					} else {
						log.Printf("Sender %s does not exist in peerMap", sender)
					}
				default:
					log.Printf("Received unknown control message from %s: Command=%s")
				}
			}
		}
		validateMessageParseDebug(parsedMsg)
		msgCount++
	}
	log.Printf("Read %d messages in total in peer: %s\n", msgCount, s.StoreOpts.ListenAddress)
}

// When broadcasting a message, we need to ensure we're encoding it properly
func (s *Store) broadcastMessage(msg p2p.Message) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(msg); err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	for _, peer := range s.PeerMap {
		if err := peer.Send(buf.Bytes()); err != nil {
			return fmt.Errorf("failed to send to peer: %w", err)
		}
	}
	return nil
}

// handleStoreFile handles writes a file with given key and broadcast it to all peers for replication
func (s *Store) handleStoreFile(key string, r io.Reader) error {
	// Copy Reader buffer
	buf := new(bytes.Buffer)
	rCopy := io.TeeReader(r, buf)

	// Store the file
	if err := s.handleFileWrite(key, rCopy); err != nil {
		return err
	}

	// Broadcast to all peers
	fromAddr, err := util.SafeStringToAddr(s.StoreOpts.ListenAddress)
	if err != nil {
		log.Fatalf("Broadcast error: %+v", err)
	}
	message := p2p.Message{
		From: fromAddr,
		Payload: p2p.DataPayload{
			Key:  key,
			Data: buf.Bytes(),
		},
	}
	return s.broadcastMessage(message)
}

// TODO: NOT BEING USED YET
func (s *Store) handlePeerControlMessages() {
	// TODO: Move control messages to different channel
	for {
		msg := <-s.Transport.Consume()
		controlMessage := p2p.MESSAGE_EXIT_CONTROL_COMMAND
		//if controlMessage == p2p.MESSAGE_EXIT_CONTROL_COMMAND {
		//	s.teardownHyperStoreServer()
		//}
		switch controlMessage {
		case p2p.MESSAGE_STORE_CONTROL_COMMAND:
			fmt.Printf("Received STORE_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		case p2p.MESSAGE_LIST_CONTROL_COMMAND:
			fmt.Printf("Received LIST_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		case p2p.MESSAGE_FETCH_CONTROL_COMMAND:
			fmt.Printf("Received FETCH_CONTROL_COMMAND from ---- %s\n", msg.From.String())
		default:
			panic("unhandled default case")
		}
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
