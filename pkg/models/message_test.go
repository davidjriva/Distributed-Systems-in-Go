package models

import (
	"encoding/json"
	"testing"
)

func TestMessage_MarshalJSON_AllFields(t *testing.T) {
	body := map[string]interface{}{
		"type": "echo",
	}

	var bodyBytes, _ = json.Marshal(body)

	msg := Message{
		Src:  "Node_1",
		Dest: "Node_2",
		Body: bodyBytes,
	}

	if got := msg.Type(); got != "echo" {
		t.Errorf("Expected Type 'echo', got %s", got)
	}
}
