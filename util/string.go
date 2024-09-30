package util

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func SafeByteToString(b []byte) (string, error) {
	if utf8.Valid(b) {
		return strings.TrimSpace(string(b)), nil
	}
	return "", fmt.Errorf("invalid UTF-8 byte sequence")
}
