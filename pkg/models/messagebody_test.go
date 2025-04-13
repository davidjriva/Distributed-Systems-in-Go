package models

import (
	"encoding/json"
	"testing"
)

func TestMessageBody_MarshalJSON_AllFields(t *testing.T) {
	body := MessageBody{
		Type:      "echo",
		MsgID:     1,
		InReplyTo: 1,
		Code:      1,
		Text:      "OK",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal MessageBody: %v", err)
	}

	expected := `{"type":"echo","msg_id":1,"in_reply_to":1,"code":1,"text":"OK"}`
	if string(data) != expected {
		t.Errorf("Unexpected JSON output. Got %s; want %s", string(data), expected)
	}
}

func TestMessageBody_MarshalJSON_RequiredFields(t *testing.T) {
	body := MessageBody{
		Type: "echo",
		Code: 1,
		Text: "OK",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal MessageBody: %v", err)
	}

	expected := `{"type":"echo","code":1,"text":"OK"}`
	if string(data) != expected {
		t.Errorf("Unexpected JSON output. Got %s; want %s", string(data), expected)
	}
}

func TestMessageBody_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"type":"error","msg_id":42,"in_reply_to":24,"code":500,"text":"Internal Error"}`

	var body MessageBody
	err := json.Unmarshal([]byte(jsonStr), &body)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if body.Type != "error" || body.MsgID != 42 || body.InReplyTo != 24 || body.Code != 500 || body.Text != "Internal Error" {
		t.Errorf("Unmarshalled MessageBody has unexpected values: %+v", body)
	}
}

func TestMessageBody_EmptyOmission(t *testing.T) {
	body := MessageBody{} // All zero-values

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal empty MessageBody: %v", err)
	}

	if string(data) != `{}` {
		t.Errorf("Expected empty JSON object for zero MessageBody, got: %s", string(data))
	}
}
