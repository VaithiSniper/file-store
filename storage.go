package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"file-store/internal/file"
	"file-store/internal/p2p"
	"file-store/internal/util"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
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
	StoreOpts              StoreOpts
	Transport              p2p.Transport
	PeerLock               sync.Mutex
	PeerMap                map[string]p2p.Peer
	FetchResponseChans     map[string]chan p2p.FetchResult
	FetchResponseChansLock sync.RWMutex
}

var globalStore *Store

// createStoreWithDefaultOptions initializes a Store with default options using a content-addressable path transform function.
func createStoreWithDefaultOptions(listenAddress string, bootstrapNodes []string, fileStorageBasePath string) *Store {
	// Prepare Transport with opts
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddress,
		HandshakeFunc: p2p.NOHANDSHAKE,
		Codec:         &p2p.DefaultCodec{},
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
		StoreOpts:              opts,
		Transport:              tTransport,
		PeerLock:               sync.Mutex{},
		PeerMap:                make(map[string]p2p.Peer),
		FetchResponseChans:     make(map[string]chan p2p.FetchResult),
		FetchResponseChansLock: sync.RWMutex{},
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
	addr, _ := util.SafeStringToAddr(s.StoreOpts.ListenAddress)
	log.Printf("Listening on %v", addr.String())

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
		senderAddr := parsedMsg.From.String()

		// Validate if peer exists
		sender, senderExists := s.PeerMap[senderAddr]
		if !senderExists {
			log.Printf("Error: Sender %s does not exist in peerMap", sender)
		}

		// Call appropriate handler
		var err error = nil
		switch parsedMsg.Type {
		case p2p.DataMessageType:
			payload := parsedMsg.Payload.(p2p.DataPayload)
			log.Printf("Parsed %s", payload.String())
			err = s.handleReadDataMessage(&payload, sender)
		case p2p.ControlMessageType:
			payload := parsedMsg.Payload.(p2p.ControlPayload)
			log.Printf("Parsed %s", payload.String())
			err = s.handleReadControlMessage(&payload, sender)
		}
		if err != nil {
			log.Printf("Error while reading message from peer %s: %v", senderAddr, err)
		}

		msgCount++
	}
	log.Printf("Read %d messages in total in peer: %s\n", msgCount, s.StoreOpts.ListenAddress)
}

func (s *Store) handleReadDataMessage(payload *p2p.DataPayload, fromPeer p2p.Peer) error {
	// If we receive a DataPayload that has fetch_id in metadata, then we are parsing a response to FETCH call
	if payload.Metadata != nil {
		fetchID, hasFetchID := payload.Metadata["fetch_id"]
		if hasFetchID {
			// So, if it is present, we push this into the corresponding fetchResponseChan
			fetchResponseChan := s.safeOperationToFetchResponseChans(util.MAP_GET_ELEMENT, fetchID, nil)
			if fetchResponseChan != nil {
				select {
				case fetchResponseChan <- p2p.FetchResult{
					FileExists: true,
					Data:       payload.Data,
					PeerAddr:   fromPeer.String(),
				}:
					log.Printf("Sent file data to waiting channel for fetch ID: %s", fetchID)
				default:
					log.Printf("Warning: Unable to send file data, channel might be full or closed for ID: %s", fetchID)
				}
				return nil
			}
		}
	}

	// If we receive a normal DataPayload, then we need to call file write for current instance
	data := bytes.NewReader(payload.Data)
	if _, err := s.handleFileWrite(payload.Key, data); err != nil {
		return err
	}
	return nil
}

func (s *Store) handleReadControlMessage(payload *p2p.ControlPayload, fromPeer p2p.Peer) error {
	/*  If we receive a ControlPayload, then we need to
	1. If Command=EXIT, then we need to remove that peer from peerMap
	2. If Command=STORE, then we need to stream a file from the sender
	*/

	log.Printf("In handleReadControlMessage with %v as COMMAND", payload.Command)

	switch payload.Command {
	case p2p.MESSAGE_EXIT_CONTROL_COMMAND:
		delete(s.PeerMap, fromPeer.String())
	case p2p.MESSAGE_STORE_CONTROL_COMMAND:
		var (
			key, keyExists          = payload.Args["key"]
			fileSizeStr, fileExists = payload.Args["size"]
		)
		if !keyExists || !fileExists {
			return fmt.Errorf("missing key/size for STORE Control Message %s", fromPeer.String())
		}

		// Store the file
		fileSize, _ := strconv.ParseInt(fileSizeStr, 10, 64)
		log.Printf("Reading streamed file of size %v", fileSize)
		_, err := s.handleFileWrite(key, io.LimitReader(fromPeer, fileSize))
		if err != nil {
			return err
		}
		// Sync wg to allow tcp read loop to continue
		fromPeer.(*p2p.TCPPeer).Wg.Done()

	case p2p.MESSAGE_LIST_CONTROL_COMMAND:
		log.Printf("Received LIST Control Message from %s", fromPeer)

	case p2p.MESSAGE_FETCH_RESPONSE_CONTROL_COMMAND:
		log.Printf("Received FETCH_RESPONSE Control Message from %s", fromPeer)
		var (
			fileFoundResp, fileFoundRespExists = payload.Args["file_exists"]
		)
		if !fileFoundRespExists {
			return fmt.Errorf("missing file_exists for FETCH_RESPONSE Control Message %s", fromPeer.String())
		}
		log.Printf("File was found on peer %s: YES/NO: %v", fromPeer.String(), fileFoundResp)

	case p2p.MESSAGE_FETCH_CONTROL_COMMAND:
		log.Printf("Received FETCH Control Message from %s", fromPeer)
		var (
			key, keyExists         = payload.Args["key"]
			fetchID, fetchIDExists = payload.Args["fetch_id"]
		)
		if !keyExists || !fetchIDExists {
			return fmt.Errorf("missing key/fetchID for FETCH Control Message %s", fromPeer.String())
		}

		// Check if file is there in this peer
		if bytesRead, err := s.handleGetFile(key, false); err != nil || bytesRead == nil {
			// Generate negative ACK and send to source
			log.Printf("File not found on this machine, sending negative ACK")
			msg := p2p.ConstructFetchResponseMessage(false)
			if err := s.broadcastMessage(msg); err != nil {
				return err
			}
		} else {
			// Generate positive ACK and send to source
			log.Printf("File found on this machine, sending ACK")
			msg := p2p.ConstructFetchResponseMessage(true)
			if err := s.sendMessageToPeer(msg, fromPeer); err != nil {
				return err
			}
			log.Printf("Sent ACK to peer %s", fromPeer.String())
			// Generate DataMessage with read file bytes and send to source
			msg = p2p.Message{
				Type: p2p.DataMessageType,
				From: nil,
				Payload: p2p.DataPayload{
					Key:  key,
					Data: bytesRead,
					Metadata: map[string]string{
						"fetch_id": fetchID,
					},
				},
			}
			log.Printf("Prepared msg: %s", msg)
			if err := s.sendMessageToPeer(msg, fromPeer); err != nil {
				return err
			}
		}
	default:
		log.Printf("Received unknown control message from %s: Command=%s", fromPeer, payload.Command)
	}

	return nil
}

// broadcastMessage broadcasts the given msg across all the peers
func (s *Store) broadcastMessage(msg p2p.Message) error {
	fromAddr, err := util.SafeStringToAddr(s.StoreOpts.ListenAddress)
	if err != nil {
		log.Fatalf("Conv error: %+v", err)
	}
	msg.From = fromAddr
	log.Printf("Broadcasting message: %+v", msg.String())
	for _, peer := range s.PeerMap {
		if err := s.Transport.(*p2p.TCPTransport).Codec.Encode(peer.(*p2p.TCPPeer).Conn, &msg); err != nil {
			return err
		}
	}
	return nil
}

// sendMessageToPeer sends the given message msg to the peer toPeer
func (s *Store) sendMessageToPeer(msg p2p.Message, toPeer p2p.Peer) error {
	fromAddr, err := util.SafeStringToAddr(s.StoreOpts.ListenAddress)
	if err != nil {
		log.Fatalf("Conv error: %+v", err)
	}
	msg.From = fromAddr
	log.Printf("Directly sending message (%s->%s): %+v", msg.From, toPeer, msg.String())
	if err := s.Transport.(*p2p.TCPTransport).Codec.Encode(toPeer.(*p2p.TCPPeer).Conn, &msg); err != nil {
		return err
	}
	return nil
}

// handleStoreFile handles writes a file with given key and broadcast it to all peers for replication
func (s *Store) handleStoreFile(key string, r io.Reader) error {
	// Copy Reader buffer
	buf := new(bytes.Buffer)
	rCopy := io.TeeReader(r, buf)

	// Store the file
	fileSize, err := s.handleFileWrite(key, rCopy)
	if err != nil {
		return err
	}

	// Now, we need to decide whether to stream	this data or to use directly send via DataPayload
	var message p2p.Message
	// If file size is beyond MaxAllowedDataPayloadSize, then decoder's buffer will overflow
	if fileSize > util.MaxAllowedDataPayloadSize {
		// Thus, we need to send a STORE control message with the necessary information to allow peers to stream
		message.Type = p2p.ControlMessageType
		message.Payload = p2p.ControlPayload{
			Command: p2p.MESSAGE_STORE_CONTROL_COMMAND,
			Args: map[string]string{
				"key":  key,
				"size": strconv.FormatInt(fileSize, 10),
			},
		}
		// Broadcast the ControlMessage
		if err := s.broadcastMessage(message); err != nil {
			return err
		}
		// And we need to stream the file contents to all peers
		for _, peer := range s.PeerMap {
			if n, err := io.Copy(peer, buf); err != nil {
				log.Printf("Streaming error: %+v", err)
				return err
			} else if n != fileSize {
				log.Printf("Streaming issue: Number of bytes streamed=%d and Number of bytes written=%d do not match", n, fileSize)
			}
		}
		log.Println("Streamed file contents to all peers successfully")
	} else {
		// Else, we can directly send a DataPayload message with the file data and key to use while replicating
		message.Type = p2p.DataMessageType
		message.Payload = p2p.DataPayload{
			Key:  key,
			Data: buf.Bytes(),
		}
		// Broadcast the ControlMessage
		if err := s.broadcastMessage(message); err != nil {
			return err
		}
	}
	return nil
}

// handleGetFile handles a file fetch with given key. If found in same store, it directly returns.
// Else broadcasts a FETCH control message to check if any peer has it.
func (s *Store) handleGetFile(key string, toBroadcast bool) ([]byte, error) {
	// TODO: Try to read and check for existence at once
	if s.existsInStorage(key) {
		bytesRead, err := s.handleFileRead(key)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Printf("File %s does not exist", key)
				return nil, os.ErrNotExist
			} else {
				return nil, err
			}
		}
		return bytesRead, nil
	}
	log.Printf("File %s does not exist in current storage, checking peers...", key)

	if toBroadcast {
		// If file is not found, need to fetch from peers
		// Using fetchID to track the FETCH request
		fetchID := s.generateFetchID(key)
		// Create a response channel to collect peer responses
		fetchResponseChan := make(chan p2p.FetchResult, len(s.PeerMap))
		// Add to map safely to track
		s.safeOperationToFetchResponseChans(util.MAP_UPSERT_ELEMENT, fetchID, fetchResponseChan)
		defer s.safeOperationToFetchResponseChans(util.MAP_DELETE_ELEMENT, fetchID, nil)
		// Prepare FETCH control msg and broadcast
		msg := p2p.Message{
			Type: p2p.ControlMessageType,
			From: nil,
			Payload: p2p.ControlPayload{
				Command: p2p.MESSAGE_FETCH_CONTROL_COMMAND,
				Args: map[string]string{
					"key":      key,
					"fetch_id": fetchID,
				},
			},
		}
		if err := s.broadcastMessage(msg); err != nil {
			return nil, err
		}

		// Wait for responses with a timeout
		timer := time.NewTimer(util.FetchMessageResponseTimeout)
		defer timer.Stop()
		// Enter read loop
		for {
			select {
			case result := <-fetchResponseChan:
				if result.Error != nil {
					log.Printf("Error from peer %s: %v", result.PeerAddr, result.Error)
					continue
				}
				if result.FileExists {
					if result.Data != nil {
						return result.Data, nil
					}
				}
			case <-timer.C:
				// Timeout reached
				return nil, fmt.Errorf("timed out waiting for fetch response")
			}
		}
	}

	return nil, nil
}

// safeOperationToFetchResponseChans thread-safely performs the action op on the s.FetchResponsesChans map based on key and value
func (s *Store) safeOperationToFetchResponseChans(op util.MAP_ACTION, key string, value chan p2p.FetchResult) chan p2p.FetchResult {
	s.FetchResponseChansLock.Lock()
	defer s.FetchResponseChansLock.Unlock()
	switch op {
	case util.MAP_GET_ELEMENT:
		return s.FetchResponseChans[key]
	case util.MAP_UPSERT_ELEMENT:
		s.FetchResponseChans[key] = value
	case util.MAP_DELETE_ELEMENT:
		delete(s.FetchResponseChans, key)
	}
	return nil
}

// --------------------------------------------------------------  END OF CONTROL PLANE --------------------------------------------------------------

// --------------------------------------------------------------  FILE HANDLING --------------------------------------------------------------

// generatePath generates and returns a path to store a file with given key
func (s *Store) generatePath(key string) string {
	hashPath := s.StoreOpts.PathTransformFunc(key, s.StoreOpts.BaseStorageLocation)
	return path.Join(s.StoreOpts.BaseStorageLocation, hashPath)
}

// generatePath generates a fetchID that can used to broadcast and keep track of a FETCH request
func (s *Store) generateFetchID(key string) string {
	fetchHash := sha1.Sum([]byte(key + "-" + s.StoreOpts.ListenAddress))
	return hex.EncodeToString(fetchHash[:])
}

// handleFileWrite writes the content from the given io.Reader to a file specified by the key within the storage system.
func (s *Store) handleFileWrite(key string, r io.Reader) (int64, error) {
	//if rc, ok := r.(io.ReadCloser); ok {
	//	defer rc.Close()
	//}

	pathname := s.generatePath(key)
	f := file.File{
		KeyPath:  key,
		BasePath: pathname,
		FileMode: util.Default,
	}
	if err := f.WriteStream(r); err != nil {
		fmt.Println("Store Error: Error occurred while writing file to storage", err)
		return 0, err
	}
	return f.FileSize, nil
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
