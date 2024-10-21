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
const FILE_KEY_PATH = "TestKeyPath"
const FILE_BASE_PATH = "TestBasePath"

// setupFile quickly sets up a File instance with provided KeyPath, BasePath and FileMode and returns it
func setupFile(t *testing.T, KeyPath string, BasePath string, FileMode os.FileMode) *File {
	var file File = File{KeyPath: KeyPath, BasePath: BasePath, FileMode: FileMode}

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
	file := setupFile(t, FILE_KEY_PATH, FILE_BASE_PATH, Default)

	assert.True(t, file.Exists())

	teardownFile(t, file, true)
}

func TestFileExistsBadInput(t *testing.T) {
	file := File{KeyPath: FILE_KEY_PATH, BasePath: FILE_BASE_PATH}
	assert.False(t, file.Exists())
}
