package util

import (
	"flag"
	"fmt"
	"strings"
)

type CommandLineArgs struct {
	ListenAddress       string
	BootstrapNodes      []string
	MetadataDBPath      string
	FileStorageBasePath string
	TestStorage         bool
}

func ParseCommandLineArgs() CommandLineArgs {
	var (
		listenAddress       string
		bootstrapNodes      string
		dbPath              string
		fileStorageBasePath string
		testStorage         bool
	)

	flag.StringVar(&listenAddress, "listen", DefaultListenAddress, "The address the hyperstore server should listen on, in <address:port> notation")
	flag.StringVar(&bootstrapNodes, "bootstrap", "", "List of bootstrapped nodes in comma separated <address:port> notation")
	flag.StringVar(&dbPath, "db", DbPath, "Path that the metadata DB will be stored in")
	flag.StringVar(&fileStorageBasePath, "file-storage-path", DefaultBaseStorageLocation, "Base path that the files will be stored in")
	flag.BoolVar(&testStorage, "test-storage", false, "Setting this to true will test the store by storing a sample file")

	flag.Parse()

	var parseListenAddress = func() string {
		// TODO: Validate if addresses are valid
		return listenAddress
	}
	var parseBootstrapNodes = func() []string {
		// TODO: Validate if addresses are valid
		if bootstrapNodes == "" {
			return make([]string, 0)
		}
		return strings.Split(bootstrapNodes, ",")
	}
	var parseDBPath = func() string {
		// TODO: Validate if path exists
		return dbPath
	}
	var parseFileStorageBasePath = func() string {
		// TODO: Validate if path exists
		if fileStorageBasePath == DefaultBaseStorageLocation {
			basePathPrefix, _ := SafeStringToAddr(listenAddress)
			defaultFileStorageBasePath := fmt.Sprintf("node-%s-%s", basePathPrefix, DefaultBaseStorageLocation)
			return defaultFileStorageBasePath
		}
		return fileStorageBasePath
	}
	var parseTestStorage = func() bool {
		// TODO: Validate if path exists
		return testStorage
	}

	flag.Parse()
	return CommandLineArgs{
		ListenAddress:       parseListenAddress(),
		BootstrapNodes:      parseBootstrapNodes(),
		MetadataDBPath:      parseDBPath(),
		FileStorageBasePath: parseFileStorageBasePath(),
		TestStorage:         parseTestStorage(),
	}
}
