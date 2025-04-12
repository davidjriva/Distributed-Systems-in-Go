package models

/*
	initmessagebody.go

	InitMessageBody represents the message body for the "init" message.
*/

type InitMessageBody struct {
	MessageBody
	NodeID  string   `json:"node_id,omitempty"`
	NodeIDs []string `json:"node_ids,omitempty"`
}
