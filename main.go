package main

import (
	"file-store/p2p"
	"fmt"
)

const hyperstoreArt = `


 | |  | \ \   / /  __ \|  ____|  __ \ / ____|__   __/ __ \|  __ \|  ____|
 | |__| |\ \_/ /| |__) | |__  | |__) | (___    | | | |  | | |__) | |__   
 |  __  | \   / |  ___/|  __| |  _  / \___ \   | | | |  | |  _  /|  __|  
 | |  | |  | |  | |    | |____| | \ \ ____) |  | | | |__| | | \ \| |____ 
 |_|  |_|  |_|  |_|    |______|_|  \_\_____/   |_|  \____/|_|  \_\______|


`

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
	fmt.Printf("Starting file-store\n%+v", hyperstoreArt)

	globalStore = getStoreInstance()
	globalStore.setupHyperStoreServer()

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
