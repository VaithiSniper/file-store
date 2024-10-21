package file

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const FILE_CONTENT = "some png bytes"
const FILE_KEY_PATH_1 = "TestKeyPath1"
const FILE_KEY_PATH_2 = "TestKeyPath2"
const FILE_BASE_PATH = "TestBasePath"

// setupFile quickly sets up a File instance with provided KeyPath, BasePath and FileMode and returns it
func setupFile(t *testing.T, KeyPath string, BasePath string, FileMode os.FileMode) *File {
	var file = File{KeyPath: KeyPath, BasePath: BasePath, FileMode: FileMode}

	data := []byte(FILE_CONTENT)
	var r io.Reader = bytes.NewReader(data)

	err := file.WriteStream(r)
	assert.Nil(t, err)

	stat, err := os.Stat(path.Join(BasePath, KeyPath))
	assert.Nil(t, err)
	assert.NotNil(t, stat)

	return &file
}

// teardownFile tears down the provided File instance
func teardownFile(t *testing.T, file *File, isErrNil bool) {
	dir := filepath.Dir(file.BasePath + "/")
	err := os.RemoveAll(dir)

	if isErrNil {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}
}

func TestFileExists(t *testing.T) {
	file := setupFile(t, FILE_KEY_PATH_1, FILE_BASE_PATH, Default)

	assert.True(t, file.Exists())

	teardownFile(t, file, true)
}

func TestFileExistsBadInput(t *testing.T) {
	file := File{KeyPath: FILE_KEY_PATH_1, BasePath: FILE_BASE_PATH}
	assert.False(t, file.Exists())
}

func TestWriteReadFile(t *testing.T) {
	file := setupFile(t, FILE_KEY_PATH_1, FILE_BASE_PATH, Default)

	readBytes, err := file.ReadFile()
	assert.Nil(t, err)
	assert.NotNil(t, readBytes)
	assert.Equal(t, string(readBytes), FILE_CONTENT)

	teardownFile(t, file, true)
}

func TestDeleteFile(t *testing.T) {
	file1 := setupFile(t, FILE_KEY_PATH_1, FILE_BASE_PATH, Default)
	file2 := setupFile(t, FILE_KEY_PATH_2, FILE_BASE_PATH, Default)

	err := file1.DeleteFile()
	assert.Nil(t, err)
	assert.False(t, file1.Exists())

	stat, err := os.Stat(path.Join(file2.BasePath, FILE_KEY_PATH_2))
	assert.Nil(t, err)
	assert.NotNil(t, stat)

	err = file2.DeleteFile()
	assert.Nil(t, err)
	assert.False(t, file1.Exists())
}
