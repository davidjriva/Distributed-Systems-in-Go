package models

import (
	"encoding/json"
	"fmt"
)

/*
	message.go

	Message represents a message sent from Src node to Dest node. The body is stored as unparsed JSON
	so the handler can parse it itself.
*/

type Message struct {
	Src  string `json:"src,omitempty"`
	Dest string `json:"dest,omitempty"`

	// A byte[] containing the serialized JSON message
	Body json.RawMessage `json:"body,omitempty"`
}

/*
Type returns the "type" field from the message body. Returns blank string if field does not exist or body is malformed.
*/
func (m *Message) Type() string {
	// The Body field is of type json.RawMessage, which is a []byte, representing a serialized JSON message.
	// We can unmarshal the message to transform it back into a MessageBody struct.
	var msgBody MessageBody

	err := json.Unmarshal([]byte(m.Body), msgBody)

	if err != nil {
		fmt.Println("Error: ", err)
		return "Error"
	}

	fmt.Printf("%+v\n", msgBody)

	return "Success"
}
