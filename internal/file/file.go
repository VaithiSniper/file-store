package file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	BasePath string
	KeyPath  string
	FileMode os.FileMode
	FileSize int64
}

// WriteStream writes into the File f from io.Reader r
func (f *File) WriteStream(r io.Reader) error {
	// Open the file and create a fd
	fd, err := f.openFileForWriting()
	if err != nil {
		log.Printf("File Error: Couldn't create file descriptor for writing: %+v", err)
		return err
	}

	writer := bufio.NewWriter(fd)
	// Copy to buffered writer
	if n, err := io.Copy(writer, r); err != nil {
		log.Printf("File Error: Error writing contents into file descriptor: %+v", err)
		return err
	} else {
		f.FileSize = n
	}
	// Flush the buffered writer
	if err := writer.Flush(); err != nil {
		log.Printf("File Error: Error flushing writer: %+v", err)
		return err
	}
	// Sync changes to disk
	if err := fd.Sync(); err != nil {
		log.Printf("File Error: Error syncing file: %+v", err)
		return err
	}

	log.Printf("Written %d bytes to %s/%s", f.FileSize, f.BasePath, f.KeyPath)
	// Close the open fd
	return fd.Close()
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
	fullPath := filepath.Join(f.BasePath, f.KeyPath)
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
