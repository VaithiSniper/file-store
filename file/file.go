package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	BasePath string
	KeyPath  string
	FileMode os.FileMode
	fileSize int64
}

// File permission octal constants
const (
	Read      os.FileMode = 0500
	Write     os.FileMode = 0200 // Change to allow write-only permission
	ReadWrite os.FileMode = 0600 // Allow both read and write for owner
	Default   os.FileMode = 0755 // Typical file permission
	All       os.FileMode = 0777 // Full permission for everyone
)

// WriteStream writes into the File f from io.Reader r
func (f *File) WriteStream(r io.Reader) error {
	fd, err := f.openFileForWriting()
	if err != nil {
		fmt.Println("File Error: Couldn't create file descriptor for writing", err)
		return err
	}
	n, err := io.Copy(fd, r)
	if err != nil {
		fmt.Println("File Error: Error writing contents into file descriptor", err)
		return err
	}
	f.fileSize = n
	return nil
}

// ReadFile reads the File f and returns byte array of content
func (f *File) ReadFile() ([]byte, error) {
	if f.Exists() {
		fullPath := fmt.Sprintf("%s/%s", f.BasePath, f.KeyPath)
		return os.ReadFile(fullPath)
	} else {
		return nil, os.ErrNotExist
	}
}

// DeleteFile deletes the File f
func (f *File) DeleteFile() error {
	if f.Exists() {
		fullPath := fmt.Sprintf("%s/%s", f.BasePath, f.KeyPath)
		if err := os.RemoveAll(fullPath); err != nil {
			return err
		}
		fmt.Println("Deleted file!")
		return f.deleteParentFolders()
	} else {
		return os.ErrNotExist
	}
}

// openFileForWriting creates the necessary subdirectories and opens a file descriptor to the File f
func (f *File) openFileForWriting() (*os.File, error) {
	if err := os.MkdirAll(f.BasePath, f.FileMode); err != nil {
		fmt.Println("File Error: Couldn't create subdirs for writing", err)
		return nil, err
	}
	fullPath := fmt.Sprintf("%s/%s", f.BasePath, f.KeyPath)
	fmt.Printf("File will be written to %s\n", fullPath)
	return os.Create(fullPath)
}

// deleteParentFolders recursively removes parent directories if they become empty
func (f *File) deleteParentFolders() error {
	dir := filepath.Dir(f.BasePath + "/")
	fmt.Printf("Deleting directories in path %s\n", dir)
	for dir != "." && dir != "/" {
		err := os.Remove(dir)
		if err != nil {
			// Stop if the directory is not empty
			if os.IsNotExist(err) || os.IsPermission(err) {
				break
			}
			return fmt.Errorf("failed to remove directory %s: %v", dir, err)
		}
		dir = filepath.Dir(dir)
	}
	return nil
}

// Exists checks if the File f Exists
func (f *File) Exists() bool {
	fullPath := fmt.Sprintf("%s/%s", f.BasePath, f.KeyPath)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println("File Error: Couldn't check if file Exists", err)
		return false
	}
	return true
}
