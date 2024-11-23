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

	// Helper funcs for testing storage
	// timeout sleeps for given seconds
	var timeout = func(seconds time.Duration) {
		time.Sleep(time.Second * seconds)
	}
	// testStoreFile tests file storing
	var testStoreFile = func(key string, useLargeFile bool) {
		stringContent := util.DefaultFileContent
		if useLargeFile {
			stringContent = util.DefaultLargeFileContent
		}
		data := bytes.NewReader([]byte(stringContent))
		if err := globalStore.handleStoreFile(key, data); err != nil {
			log.Fatalf("Error while writing test file -> %+v", err)
		}
	}
	// testGetFile tests file retrieval
	var testGetFile = func(key string) {
		if bytesRead, err := globalStore.handleGetFile(key, true); err != nil {
			log.Fatalf("Error while getting test file -> %+v", err)
		} else {
			log.Printf("Successfully got test file contents -> %s", string(bytesRead))
		}
	}
	// testDeleteFile deletes the file locally
	var testDeleteFile = func(key string) {
		if err := globalStore.handleFileDelete(key); err != nil {
			log.Fatalf("Error while deleting test file -> %+v", err)
		}
	}

	// Test out storage functionality
	if commandLineArgs.TestStorage {
		log.Println("Basic storage FT")
		timeout(2)
		testStoreFile("test_key", false)
		timeout(2)
		testGetFile("test_key")
		timeout(2)
		testDeleteFile("test_key")
		timeout(5)
		testGetFile("test_key")
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
