package main

import (
	"file-store/internal/db"
	"file-store/internal/p2p"
	"fmt"
	"log"
)

const hyperstoreArt = `


 | |  | \ \   / /  __ \|  ____|  __ \ / ____|__   __/ __ \|  __ \|  ____|
 | |__| |\ \_/ /| |__) | |__  | |__) | (___    | | | |  | | |__) | |__   
 |  __  | \   / |  ___/|  __| |  _  / \___ \   | | | |  | |  _  /|  __|  
 | |  | |  | |  | |    | |____| | \ \ ____) |  | | | |__| | | \ \| |____ 
 |_|  |_|  |_|  |_|    |______|_|  \_\_____/   |_|  \____/|_|  \_\______|


`

func onPeerFailure(peer p2p.Peer) error {
	return fmt.Errorf("error occuring")
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

	ddbInstance, err := db.InitDB(db.DB_PATH)
	if err != nil {
		log.Fatalf("error occurred while setting up ddb: %+v\n", err)
	}

	if ddbInstance.IsInit {
		fmt.Println("ddb instance is initialized")
	}
	if ddbInstance.IsReady {
		fmt.Println("ddb instance is ready for tx")
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
