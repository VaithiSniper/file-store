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
	return peer.Close()
}

func main() {
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
