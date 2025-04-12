package models

/*
	messagebody.go

	MessageBody represents the reserved keys for a message body.
*/

type MessageBody struct {
	// Message type.
	//Could be request, response, or error.
	// omitempty; If this field is empty, then omit it from the serialized JSON output.
	Type string `json:"type,omitempty"`

	// Optional. Message identifier that is unique to the source node.
	MsgID int `json:"msg_id,omitempty"`

	// Optional. For request/response, the msg_id of the request.
	InReplyTo int `json:"in_reply_to,omitempty"`

	// Error code, if an error occurred.
	Code int `json:"code,omitempty"`

	// Error message, if an error occurred.
	Text string `json:"text,omitempty"`
}
