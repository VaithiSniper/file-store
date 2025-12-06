package main

import (
	"bytes"
	"file-store/internal/db"
	"file-store/internal/logger"
	"file-store/internal/util"
)

func initApp() {
	util.RegisterGobTypes()
}

func basicStoreSmokeTest() {
	const moduleName = "SMOKE_TEST"

	// testStoreFile tests file storing
	var testStoreFile = func(key string, useLargeFile bool) {
		stringContent := util.DefaultFileContent
		if useLargeFile {
			stringContent = util.DefaultLargeFileContent
		}
		data := bytes.NewReader([]byte(stringContent))
		if err := globalStore.handleStoreFile(key, data); err != nil {
			logger.LogEmergency(moduleName, "Error while writing test file: %+v", err)
		}
	}
	// testGetFile tests file retrieval
	var testGetFile = func(key string) {
		if bytesRead, err := globalStore.handleGetFile(key, true); err != nil {
			if err.Error() == "Timed out waiting for fetch response." {
				logger.LogWarning(
					moduleName,
					"Couldn't find files in peers. Maybe nodes are not bootstrapped?",
				)
			} else {
				logger.LogEmergency(
					moduleName, "Error while getting test file: %+v", err,
				)
			}
		} else {
			logger.LogNotice(
				"Successfully got test file contents -> %s", string(bytesRead),
			)
		}
	}
	// testDeleteFile deletes the file locally
	var testDeleteFile = func(key string) {
		if err := globalStore.handleFileDelete(key); err != nil {
			logger.LogEmergency(
				moduleName, "Error while deleting test file -> %+v", err,
			)
		}
	}

	// Test out storage functionality
	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file storing")
	testStoreFile("test_key", false)
	logger.LogForce(moduleName, "FILE STORAGE: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file retrieval")
	testGetFile("test_key")
	logger.LogForce(moduleName, "FILE RETRIEVAL: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing file deletion")
	testDeleteFile("test_key")
	logger.LogForce(moduleName, "FILE DELETION: PASSED")

	util.TimeoutBySeconds(2)
	util.PrintInBanner("Testing retrieval of deleted file (should error)")
	testGetFile("test_key")
	logger.LogForce(moduleName, "FILE GET AFTER DELETE: PASSED")

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
	moduleName := "INIT_DDB"
	ddbInstance, err := db.InitDB(util.DbPath)
	if err != nil {
		logger.LogEmergency(
			moduleName,
			"error occurred while setting up ddb: %+v\n", err,
		)
	}

	if ddbInstance.IsInit {
		logger.LogNotice(moduleName, "ddb instance is initialized")
	}
	if ddbInstance.IsReady {
		logger.LogNotice(moduleName, "ddb instance is ready for tx")
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
