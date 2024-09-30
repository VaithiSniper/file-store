package main

import (
	"file-store/p2p"
	"file-store/util"
	"fmt"
	"log"
	"sync"
)

// TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.
func onPeerFailure(peer p2p.Peer) error {
	return fmt.Errorf("Error occuring")
}

func onPeerSuccess(peer p2p.Peer) error {
	return nil
}

func onPeerAbruptPeerCloseFailure(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	fmt.Println("Starting file-store")

	var wg sync.WaitGroup

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":5000",
		HandshakeFunc: p2p.NOHANDSHAKE,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        onPeerSuccess,
	}

	tTransport := p2p.NewTCPTransport(tcpOpts)

	globalStore = getStoreInstance()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg := <-tTransport.Consume()
			fmt.Printf("%+v\n", msg.Payload)
			str, err := util.SafeByteToString(msg.Payload)
			if err != nil {
				log.Println("Error converting bytes to string:", err)
				continue // Instead of panicking, just log the error and continue
			}
			fmt.Printf("%+v\n", str)
			if str == "" {
				continue
			}
			// controlMessageType := globalStore.parseControlMessageType(str)
		}
	}()

	fmt.Println("Starting to listen and accept connections...")
	if err := tTransport.ListenAndAccept(); err != nil {
		log.Fatalln("Error listening and accepting connections:", err)
	}

	// Wait for the goroutine to finish (which would be never in this case)
	wg.Wait()

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
