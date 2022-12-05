package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/yomorun/yomo/core/frame"
)

// Message is the structure that we received from YoMo Zipper.
type Message struct {
	Meta *MessageMeta `json:"meta"`
	Data []byte       `json:"data"`
}

// MessageMeta describes the meta data of a message.
type MessageMeta struct {
	Mode     string `json:"mode"` // datagram | stream
	StreamID int64  `json:"stream_id"`
}

// DataTags specify the tags that we observed from zipper.
func DataTags() []byte {
	return []byte{0x30}
}

// Handler is the function that will be called when we receive a message.
// No matter payload is over `Datagram` or `Stream`, Handler will be called when we receive a message.
func Handler(tag frame.Tag, request []byte) (byte, []byte) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	var message Message
	if err := json.Unmarshal(request, &message); err != nil {
		return 0, nil
	}

	// Convert to uppercase
	log.Printf("[mode=%s][streamID=%d] received data, len=%d, data:%+v", message.Meta.Mode, message.Meta.StreamID, len(message.Data), message.Data)
	message.Data = []byte(strings.ToUpper(string(message.Data)))

	// Send data back to webtransport.day client.
	response, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return 0x31, response
}
