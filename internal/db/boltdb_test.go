package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
	"testing"
	"time"
)

// setupDB quickly sets up a DDB instance with provided dbPath and returns it
func setupDB(t *testing.T, dbPath string) *DDB {
	ddb, err := InitDB(dbPath)
	assert.Nil(t, err)
	assert.NotNil(t, ddb)

	return &ddb
}

// teardownDB tears down any existing singleton instance of DDB
func teardownDB(t *testing.T, isErrNil bool) {
	CloseDB()
	err := TeardownDB(DB_PATH)
	if isErrNil {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}
}

// getSampleKeyValuePairs generates and populates a sample 5 membered key-value pair map that can be used for tests
func getSampleKeyValuePairs() map[string]string {
	numbers := [5]int{11, 22, 33, 44, 55}
	var sampleMap = make(map[string]string)
	for n := range numbers {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		sampleMap[key] = value
	}
	return sampleMap
}

// --------------------------------------------------------------  DB MANAGEMENT TESTS --------------------------------------------------------------

func TestInitDB(t *testing.T) {
	ddb := setupDB(t, DB_PATH)

	// Check members
	assert.NotNil(t, ddb.db)
	assert.Equal(t, ddb.dbPath, DB_PATH)
	assert.True(t, ddb.IsInit)
	assert.True(t, ddb.IsReady)
	assert.NotNil(t, ddb.CreatedAt)
	assert.LessOrEqual(t, time.Since(ddb.CreatedAt), time.Millisecond*100)

	// Check bbolt instance
	assert.Equal(t, ddb.db.Path(), DB_PATH)
	err := ddb.db.View(func(tx *bbolt.Tx) error {
		var bucketNameBytes = []byte(METADATA_BUCKET_NAME)
		b := tx.Bucket(bucketNameBytes)
		assert.NotNil(t, b)
		return nil
	})
	assert.Nil(t, err)

	// Teardown
	teardownDB(t, true)
}

func TestTeardownDBWithoutInit(t *testing.T) {
	teardownDB(t, false)
}

func TestTeardownDBAfterInit(t *testing.T) {
	_ = setupDB(t, DB_PATH)
	teardownDB(t, true)
}

func TestGetDB(t *testing.T) {
	ddb := setupDB(t, DB_PATH)

	assert.NotNil(t, ddb.db)
	assert.Equal(t, ddb, GetDB())
}

func TestCloseDB(t *testing.T) {
	ddb := setupDB(t, DB_PATH)

	assert.NotNil(t, ddb.db)
	CloseDB()
	assert.Empty(t, ddb.db.Path())
}

// --------------------------------------------------------------  DB MANAGEMENT TESTS --------------------------------------------------------------

// --------------------------------------------------------------  DB CRUD TESTS --------------------------------------------------------------

func TestSetGetValue(t *testing.T) {
	ddb := setupDB(t, DB_PATH)
	assert.NotNil(t, ddb.db)

	sampleKeyValuePairs := getSampleKeyValuePairs()
	for k, v := range sampleKeyValuePairs {
		ddb.setValue(k, v)
	}

	for k, v := range sampleKeyValuePairs {
		assert.Equal(t, ddb.getValue(k), v)
	}

	teardownDB(t, true)
}

// --------------------------------------------------------------  DB CRUD TESTS --------------------------------------------------------------