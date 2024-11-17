package main

import (
	"bytes"
	"file-store/internal/db"
	"file-store/internal/util"
	"fmt"
	"log"
	"time"
)

func initApp() {
	util.RegisterGobTypes()
}

func initStore(commandLineArgs util.CommandLineArgs) {
	globalStore = getStoreInstance(commandLineArgs.ListenAddress, commandLineArgs.BootstrapNodes, commandLineArgs.FileStorageBasePath)
	go globalStore.setupHyperStoreServer()

	// Test out storage functionality
	if commandLineArgs.TestStorage {
		log.Println("Basic storage FT")
		time.Sleep(time.Second * 2)
		data := bytes.NewReader([]byte(util.DefaultLargeFileContent))
		err := globalStore.handleStoreFile("test_key", data)
		if err != nil {
			log.Fatalf("Error while writing test file -> %+v", err)
		}
	}

}

func initDDB() {
	ddbInstance, err := db.InitDB(util.DbPath)
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

func keepAlive() {
	for {

	}
}

func main() {
	initApp()

	commandLineArgs := util.ParseCommandLineArgs()

	util.ColorPrint(util.ColorBlue, util.HyperstoreArt)
	log.Println("Starting file-store...")

	initStore(commandLineArgs)

	// initDDB()

	keepAlive()
}
