package storage

import (
	"bytes"
	"file-store/internal/constants"
	"file-store/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// getHashPath is a helper function to get the generated hash portion of the path
func getHashPath(fullpath string, baseStorageLocation string) string {
	if filepath.IsLocal(baseStorageLocation) {
		absPath, err := filepath.Abs(baseStorageLocation)
		if err != nil {
			fmt.Println(err)
		}
		baseStorageLocation = filepath.Base(absPath)
	}
	hashedPart := strings.TrimPrefix(fullpath, baseStorageLocation+"/")
	return hashedPart
}

func setupTestStore() *Store {
	storeOpts := StoreOpts{
		ListenAddress:       ":5000",
		BootstrapNodes:      []string{""},
		BaseStorageLocation: "",
		PathTransformFunc:   ContentAddressableTransformFunc,
	}
	store := CreateStoreWithUserOptions(storeOpts)
	return store
}

func TestContentAddressableTransformFunc(t *testing.T) {
	store := setupTestStore()

	pathOutput := store.generatePath(constants.CommonFileKey)
	hashOutput := getHashPath(pathOutput, store.StoreOpts.BaseStorageLocation)
	fmt.Println(hashOutput)
	assert.NotEqual(t, hashOutput, constants.CommonFileKey)

	regexPattern := `^([a-f0-9]{10}/){3}[a-f0-9]{10}$`
	match, err := regexp.MatchString(regexPattern, hashOutput)
	assert.Nil(t, err)
	assert.True(t, match)
}

func TestUploadFile(t *testing.T) {
	store := setupTestStore()
	data := []byte(constants.CommonStringContent)
	fileSize, err := store.handleFileWrite(
		constants.CommonFileKey, bytes.NewReader(data),
	)
	assert.Nil(t, err)
	assert.NotZero(t, fileSize)
}

func TestReadFile(t *testing.T) {
	store := setupTestStore()
	content, err := store.handleFileRead(constants.CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	sContent, err := util.SafeByteToString(content)
	fmt.Printf("Content read: %s\n", sContent)
	assert.Nil(t, err)
	assert.Equal(t, constants.CommonStringContent, sContent)
}

func TestDeleteFile(t *testing.T) {
	store := setupTestStore()
	err := store.handleFileDelete(constants.CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	assert.Nil(t, err)
}
