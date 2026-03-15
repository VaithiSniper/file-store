package file

import (
	"bytes"
	"file-store/internal/constants"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupFile quickly sets up a File instance with provided KeyPath, BasePath and FileMode and returns it
func setupFile(
	t *testing.T, KeyPath string, BasePath string, FileMode os.FileMode,
) *File {
	var file = File{KeyPath: KeyPath, BasePath: BasePath, FileMode: FileMode}

	data := []byte(constants.DefaultFileContent)
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
	file := setupFile(
		t, constants.DefaultFileKeyPath, constants.DefaultFileBasePath, constants.Default,
	)

	assert.True(t, file.Exists())

	t.Cleanup(
		func() {
			teardownFile(t, file, true)
		},
	)
}

func TestFileExistsBadInput(t *testing.T) {
	file := File{
		KeyPath: constants.DefaultFileKeyPath, BasePath: constants.DefaultFileBasePath,
	}
	assert.False(t, file.Exists())
}
