package main

import (
	"bytes"
	"file-store/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	CommonStringContent string = "some png bytes"
	CommonFileKey       string = "testfilename"
)

func TestDefaultTransformFunc(t *testing.T) {
	store := getStoreInstance()
	assert.Equal(t, store.StoreOpts.PathTransformFunc(CommonFileKey), CommonFileKey)
}

//func TestControlCommands(t *testing.T) {
//	store := createStoreWithDefaultOptions()
//}

func TestUploadFile(t *testing.T) {
	store := getStoreInstance()
	data := []byte(CommonStringContent)
	err := store.handleFileWrite(CommonFileKey, bytes.NewReader(data))
	assert.Nil(t, err)
}

func TestReadFile(t *testing.T) {
	store := getStoreInstance()
	content, err := store.handleFileRead(CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	sContent, err := util.SafeByteToString(content)
	assert.Nil(t, err)
	assert.Equal(t, CommonStringContent, sContent)
}

func TestDeleteFile(t *testing.T) {
	store := getStoreInstance()
	err := store.deleteFileFromStorage(CommonFileKey)
	// No errors should occur except file not found error
	if err != nil {
		assert.True(t, os.IsNotExist(err))
		return
	}
	assert.Nil(t, err)
}
