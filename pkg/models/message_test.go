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
		t.Errorf("Expected Type 'echo', got '%s'", got)
	}
}

func TestMessage_RPCError_ValidError(t *testing.T) {
	body := map[string]interface{}{
		"type": "echo",
		"text": "An error has occured :(",
		"code": 13,
	}

	var bodyBytes, _ = json.Marshal(body)

	msg := Message{
		Src:  "Node_1",
		Dest: "Node_2",
		Body: bodyBytes,
	}

	var rpcErr *RPCError = msg.RPCError()

	expectedErrStr := `RPCError(Crash, "An error has occured :(")`
	actualErrStr := rpcErr.Error()
	if actualErrStr != expectedErrStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedErrStr, actualErrStr)
	}
}