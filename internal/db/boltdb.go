package db

import (
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const DB_PATH = "./data/metadata.db"
const METADATA_BUCKET_NAME = "fileMetadata"

type DDB struct {
	db        *bbolt.DB
	dbPath    string
	IsInit    bool
	IsReady   bool
	CreatedAt time.Time
}

var (
	ddbInstance DDB
)

// --------------------------------------------------------------  DB MANAGEMENT --------------------------------------------------------------

// InitDB initializes the database at the specified path and creates the necessary buckets.
func InitDB(dbPath string) (DDB, error) {

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return DDB{}, fmt.Errorf("failed to create directory: %v", err)
	}

	_db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return ddbInstance, fmt.Errorf("failed to open BoltDB: %v", err)
	}

	// Create required buckets
	err = _db.Update(func(tx *bbolt.Tx) error {
		b := getBucketInstance(tx, METADATA_BUCKET_NAME)
		if b == nil {
			return fmt.Errorf("could not create bucket with name: %s", METADATA_BUCKET_NAME)
		}
		return nil
	})

	if err != nil {
		return ddbInstance, fmt.Errorf("failed to create buckets: %v", err)
	}

	ddbInstance = DDB{db: _db, dbPath: dbPath, IsInit: true, IsReady: true, CreatedAt: time.Now()}
	return ddbInstance, nil
}

// TeardownDB tears down the embedded db's files
func TeardownDB(dbPath string) error {
	// Check if the file or directory exists
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("database file or directory does not exist: %v", dbPath)
	} else if err != nil {
		return fmt.Errorf("error checking database file or directory: %v", err)
	}

	// Get the directory path
	dir := filepath.Dir(dbPath)

	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to remove database files: %v", err)
	}
	return nil
}

// GetDB returns a pointer to the existing DDB instance.
func GetDB() *DDB {
	return &ddbInstance
}

// CloseDB closes the db connection
func CloseDB() {
	if ddbInstance.IsInit && ddbInstance.IsReady {
		err := ddbInstance.db.Close()
		if err != nil {
			log.Fatalf("error closing database: %+v\n", err)
		}
		log.Println("Database connection closed")
	} else {
		log.Println("Database connection is already closed!")
	}
}

// -------------------------------------------------------------- END OF DB MANAGEMENT --------------------------------------------------------------

// --------------------------------------------------------------  DB CRUD --------------------------------------------------------------

// TODO: Implement batch reads and batch writes

// getValue gets the value for the given key string
func (ddb *DDB) getValue(key string) string {
	var valueBytes []byte
	err := ddb.db.View(func(tx *bbolt.Tx) error {
		keyBytes := []byte(key)
		b := getBucketInstance(tx, METADATA_BUCKET_NAME)
		valueBytes = b.Get(keyBytes)
		return nil
	})
	if err != nil {
		log.Fatalf("error reading key from database: %+v\n", err)
	}
	return string(valueBytes)
}

// setValue PUTS the given key-value pair in the db, i.e, overwrites if already exists
func (ddb *DDB) setValue(key string, value string) {
	err := ddb.db.Update(func(tx *bbolt.Tx) error {
		keyBytes := []byte(key)
		valueBytes := []byte(value)
		b := getBucketInstance(tx, METADATA_BUCKET_NAME)
		return b.Put(keyBytes, valueBytes)
	})
	if err != nil {
		log.Printf("Failed to set key %s in metadata bucket\n", string(key))
	}
}

// getBucketInstance returns an existing bucket or creates a new one if it doesn't exist.
func getBucketInstance(tx *bbolt.Tx, bucketName string) *bbolt.Bucket {
	bName := []byte(bucketName)
	b := tx.Bucket(bName)
	if b != nil {
		// TODO: Convert below to debug log after setting up logger
		// fmt.Printf("Found bucket with name %s\n", bucketName)
		return b
	}
	b, err := tx.CreateBucket(bName)
	if err != nil {
		fmt.Printf("Failed to create bucket %s due to error: %+v\n", bucketName, err)
	}
	return b
}

// -------------------------------------------------------------- END OF DB CRUD --------------------------------------------------------------
