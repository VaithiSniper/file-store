package main

import (
	"bytes"
	"file-store/internal/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
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

func TestContentAddressableTransformFunc(t *testing.T) {
	store := getStoreInstance()

	pathOutput := store.generatePath(util.CommonFileKey)
	hashOutput := getHashPath(pathOutput, store.StoreOpts.BaseStorageLocation)
	fmt.Println(hashOutput)
	assert.NotEqual(t, hashOutput, util.CommonFileKey)

	regexPattern := `^([a-f0-9]{10}/){3}[a-f0-9]{10}$`
	match, err := regexp.MatchString(regexPattern, hashOutput)
	assert.Nil(t, err)
	assert.True(t, match)
}

func TestUploadFile(t *testing.T) {
	store := getStoreInstance()
	data := []byte(util.CommonStringContent)
	err := store.handleFileWrite(util.CommonFileKey, bytes.NewReader(data))
	assert.Nil(t, err)
}

func TestReadFile(t *testing.T) {
	store := getStoreInstance()
	content, err := store.handleFileRead(util.CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	sContent, err := util.SafeByteToString(content)
	fmt.Printf("Content read: %s\n", sContent)
	assert.Nil(t, err)
	assert.Equal(t, util.CommonStringContent, sContent)
}

func TestDeleteFile(t *testing.T) {
	store := getStoreInstance()
	err := store.handleFileDelete(util.CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	assert.Nil(t, err)
}
