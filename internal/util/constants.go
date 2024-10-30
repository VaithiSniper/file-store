package util

import (
	"os"
)

// --------------------------------------------------------------  FILE CONSTANTS --------------------------------------------------------------

// Default File instance opts
const (
	DefaultFileContent  = "some png bytes"
	DefaultFileKeyPath  = "TestKeyPath"
	DefaultFileBasePath = "TestBasePath"
)

// File permission octal constants
const (
	Read      os.FileMode = 0500
	Write     os.FileMode = 0200 // Change to allow write-only permission
	ReadWrite os.FileMode = 0600 // Allow both read and write for owner
	Default   os.FileMode = 0755 // Typical file permission
	All       os.FileMode = 0777 // Full permission for everyone
)

// Default File content and name
const (
	CommonStringContent string = "some png bytes"
	CommonFileKey       string = "testfilename"
)

// --------------------------------------------------------------  END OF FILE CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  STORAGE CONSTANTS --------------------------------------------------------------

const (
	DefaultBaseStorageLocation string = "./storage"
	DefaultListenAddress       string = ":5000"
)

// --------------------------------------------------------------  END OF STORAGE CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  DB CONSTANTS --------------------------------------------------------------

const DbPath = "./data/metadata.db"
const MetadataBucketName = "fileMetadata"

// --------------------------------------------------------------  END OF DB CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  P2P CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  END OF P2P CONSTANTS --------------------------------------------------------------
