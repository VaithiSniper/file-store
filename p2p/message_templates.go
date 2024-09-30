package p2p

import (
	"bytes"
	"encoding/json"
	"file-store/file"
	"fmt"
	"text/template"
)

type MessageWrapper struct {
	ControlMessage ControlMessage
	File           file.File
}

type MessageFormat interface {
	generateMessage(messageWrapper *MessageWrapper) string
	parseMessage(messageString string) *MessageWrapper
}

type MessageFormatOpts struct {
	MessageFormatter MessageFormat
}

type JSONFormat struct{}

type ProtoFormat struct{}

type MessageFormatFactory struct {
	MessageFormatOpts
}

func NewMessageFormatFactory(opts MessageFormatOpts) *MessageFormatFactory {
	return &MessageFormatFactory{
		MessageFormatOpts: opts,
	}
}

func (msgFormat JSONFormat) parseMessage(messageString string) *MessageWrapper {
	return &MessageWrapper{}
}

func (msgFormat ProtoFormat) parseMessage(messageString string) *MessageWrapper {
	return &MessageWrapper{}
}

func (msgFormat JSONFormat) generateMessage(messageWrapper *MessageWrapper) string {
	jsonMessage := jsonMessageTemplate{
		ControlMessageType: messageWrapper.ControlMessage.String(),
		File:               messageWrapper.File,
	}
	marshal, err := json.Marshal(jsonMessage)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return string(marshal)
}

func (msgFormat ProtoFormat) generateMessage(messageWrapper *MessageWrapper) string {
	tmpl, err := template.New("controlMessage").Parse(protoMessageTemplate)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, tmpl)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

const protoMessageTemplate = `{
{{ @type: .ControlMessage }},
{{ if .File }}
  file:  {
	  {{ .File.KeyPath }},
	  {{ .FileKey.BasePath }},
  }
{{ end }}
}
`

type jsonMessageTemplate struct {
	ControlMessageType string    `json:"@type"`
	File               file.File `json:"file,omitempty"`
}
