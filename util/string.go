package util

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// SafeByteToString converts a byte slice to a trimmed UTF-8 string. Returns an error if the byte slice is not valid UTF-8.
func SafeByteToString(b []byte) (string, error) {
	if utf8.Valid(b) {
		return strings.TrimSpace(string(b)), nil
	}
	return "", fmt.Errorf("invalid UTF-8 byte sequence")
}

// ChunkString chunks the given hex string s into blocks of fixed size uint8 blockSize
func ChunkString(s string, blockSize uint8) []string {
	var chunks []string
	i := uint8(0)
	strLen := uint8(len(s))
	for i = 0; i < strLen; i += blockSize {
		end := i + blockSize
		if end > strLen {
			end = strLen
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
