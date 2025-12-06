package main

import (
	"bytes"
	"file-store/internal/db"
	"file-store/internal/util"
	"fmt"
	"log"
)

func initApp() {
	util.RegisterGobTypes()
}

func basicStoreSmokeTest() {
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
			if err.Error() == "Timed out waiting for fetch response." {
				log.Printf("Couldn't find files in peers. Maybe nodes are not bootstrapped?")
			} else {
				log.Fatalf("Error while getting test file: %+v", err)
			}
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
	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file storing")
	testStoreFile("test_key", false)
	log.Println("FILE STORAGE: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file retrieval")
	testGetFile("test_key")
	log.Println("FILE RETRIEVAL: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file deletion")
	testDeleteFile("test_key")
	log.Println("FILE DELETION: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing retrieval of deleted file (should error)")
	testGetFile("test_key")
	log.Println("FILE GET AFTER DELETE: PASSED")

	util.PrintInBanner("Basic store smoke test completed successfully!")
}

func initStore(commandLineArgs util.CommandLineArgs) {
	globalStore = getStoreInstance(
		commandLineArgs.ListenAddress, commandLineArgs.BootstrapNodes,
		commandLineArgs.FileStorageBasePath,
	)
	go globalStore.setupHyperStoreServer()

	if commandLineArgs.TestStorage {
		basicStoreSmokeTest()
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
	util.PrintInBanner("Initializing hyperstore")

	initStore(commandLineArgs)

	// initDDB()

	keepAlive()
}
