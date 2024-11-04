package util

import (
	"os"
)

// --------------------------------------------------------------  FILE CONSTANTS --------------------------------------------------------------

// Default File instance opts
const (
	DefaultFileContent      = "some png bytes"
	DefaultLargeFileContent = "\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit. In bibendum ipsum et arcu placerat, sit amet scelerisque dui tempor. Nunc ligula nulla, faucibus vulputate laoreet vel, ultricies nec justo. Sed volutpat lorem enim, vel imperdiet velit fringilla non. Donec nec nibh libero. Integer non dui dictum lorem semper interdum. Duis et finibus leo. Curabitur pellentesque, erat sed vestibulum elementum, augue ex vehicula quam, ac dictum lectus sem vel odio. In vitae pharetra est. Vestibulum venenatis tortor et placerat aliquet. Aenean interdum nisl at pulvinar congue. Sed odio diam, aliquam non urna ac, commodo sollicitudin elit. Sed eu purus ultrices tortor iaculis vestibulum. Vivamus vitae mi laoreet nulla faucibus pharetra. Mauris posuere sit amet dolor quis eleifend. Nunc odio nisi, accumsan non justo non, volutpat vulputate nisi. Donec sapien nisi, molestie quis augue ac, condimentum rutrum orci.\n\nEtiam condimentum, orci a euismod maximus, tortor ipsum scelerisque leo, vitae condimentum arcu risus at sem. In ut tellus eros. In elit nunc, volutpat auctor sollicitudin quis, pretium eu est. Vivamus rhoncus hendrerit neque, et ultrices mauris. Cras sollicitudin, elit sit amet interdum efficitur, arcu massa porta turpis, nec efficitur tellus nisl sed neque. Nulla vel elit felis. Maecenas eu tellus ut enim porttitor laoreet vitae non ligula.\n\nPellentesque nulla eros, finibus a nunc id, varius porta erat. Integer vestibulum eleifend ipsum in pharetra. Suspendisse potenti. Nulla ut dui libero. Phasellus scelerisque vestibulum ex vel lacinia. Proin nec tellus at orci feugiat facilisis eget eu est. Mauris sed enim ac lorem eleifend ullamcorper. Fusce vel tortor mattis massa mattis pretium.\n\nNulla feugiat vulputate leo. Nunc condimentum nibh vitae lorem placerat commodo. Cras luctus libero nec dolor luctus ornare. Quisque vitae nisi id libero bibendum convallis. Ut sit amet augue risus. Praesent at pretium purus. Pellentesque tincidunt non mi at posuere. Vivamus ullamcorper, nulla a lacinia sodales, dui leo bibendum augue justo.\n"
	DefaultFileKeyPath      = "TestKeyPath"
	DefaultFileBasePath     = "TestBasePath"
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
	DefaultBaseStorageLocation string = "storage"
	DefaultListenAddress       string = ":5000"
)

// --------------------------------------------------------------  END OF STORAGE CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  DB CONSTANTS --------------------------------------------------------------

const (
	DbPath             = "./data/metadata.db"
	MetadataBucketName = "fileMetadata"
)

// --------------------------------------------------------------  END OF DB CONSTANTS --------------------------------------------------------------

// --------------------------------------------------------------  P2P CONSTANTS --------------------------------------------------------------

const (
	DefaultChunkSize          uint8 = 10
	MessageChanBufferSize           = 32
	MaxAllowedDataPayloadSize       = 1024
)

// --------------------------------------------------------------  END OF P2P CONSTANTS --------------------------------------------------------------
