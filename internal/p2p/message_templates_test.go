package p2p

import (
	"file-store/internal/file"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEmptyMessageWrapper() MessageWrapper {
	return MessageWrapper{
		ControlMessage: MESSAGE_LIST_CONTROL_COMMAND,
		File:           file.File{KeyPath: "", BasePath: ""},
	}
}

func TestGenerateJSON(t *testing.T) {
	messageFormatFactory := NewMessageFormatFactory(MessageFormatOpts{
		MessageFormatter: JSONFormat{},
	})
	msg := getEmptyMessageWrapper()
	str := messageFormatFactory.MessageFormatter.generateMessage(&msg)

	fmt.Println(str)

	assert.NotNil(t, str)
	assert.NotEmpty(t, str)

	assert.Contains(t, str, MESSAGE_LIST_CONTROL_COMMAND.String())
	assert.Contains(t, str, msg.File.KeyPath)
	assert.Contains(t, str, msg.File.BasePath)
}

func TestGenerateProto(t *testing.T) {
	assert.True(t, true)
}
