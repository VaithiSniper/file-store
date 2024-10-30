package main

import (
	"file-store/internal/db"
	"file-store/internal/util"
	"fmt"
	"log"
)

func main() {

	commandLineArgs := util.ParseCommandLineArgs()

	util.ColorPrint(util.ColorBlue, util.HyperstoreArt)
	log.Println("Starting file-store...")

	globalStore = getStoreInstance(commandLineArgs.ListenAddress, commandLineArgs.BootstrapNodes)
	globalStore.setupHyperStoreServer()

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
