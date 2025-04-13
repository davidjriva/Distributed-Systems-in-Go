package models

import (
	"errors"
	"testing"
	"encoding/json"
	"fmt"
)

func TestRPCError_NewRPCError(t *testing.T) {
	var rpcErr *RPCError = NewRPCError(10, "This operation is not supported!")

	if rpcErr.Code != 10 {
		t.Errorf("Expected Code '10', got %d", rpcErr.Code)
	}

	if rpcErr.Text != "This operation is not supported!" {
		t.Errorf("Expected 'This operation is not supported!', got %s", rpcErr.Text)
	}
}

func TestRPCError_ErrorCodeText_AllFields(t *testing.T) {
	errCodes := map[int]string{
		0:  "Timeout",
		10: "NotSupported",
		11: "TemporarilyUnavailable",
		12: "MalformedRequest",
		13: "Crash",
		14: "Abort",
		20: "KeyDoesNotExist",
		21: "KeyAlreadyExists",
		22: "PreconditionFailed",
		30: "TxnConflict",
	}

	for code := range errCodes {
		var actualCodeText string = ErrorCodeText(code)
		if actualCodeText != errCodes[code] {
			t.Errorf("Expected '%s', got '%s'", errCodes[code], actualCodeText)
		}
	}
}

func TestRPCError_ErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "RPCError with Timeout code",
			err:      NewRPCError(Timeout, "Timeout occurred"),
			expected: Timeout,
		},
		{
			name:     "RPCError with NotSupported code",
			err:      NewRPCError(NotSupported, "Not supported"),
			expected: NotSupported,
		},
		{
			name:     "Non-RPCError",
			err:      errors.New("a regular error"),
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorCode(tt.err)
			if got != tt.expected {
				t.Errorf("ErrorCode() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestRPCError_ErrorString(t *testing.T) {
	var rpcErr *RPCError = NewRPCError(10, "This operation is not supported!")

	expectedErrorStr := "RPCError(NotSupported, \"This operation is not supported!\")"
	if rpcErr.Error() != expectedErrorStr{
		t.Errorf("Expected '%s', got '%s'", expectedErrorStr, rpcErr.Error())
	}
}

func TestRPCError_MarshalJSON(t *testing.T) {
	err := NewRPCError(10, "This operation is not supported!")

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Unexpected marshal error: %v", marshalErr)
	}

	expectedJSON := `{"type":"error","code":10,"text":"This operation is not supported!"}`
	if string(data) != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, data)
	}
}

/*
	=========================================
	Tests from the official maelstrom package
	=========================================
*/

func TestErrorCodeTextOther(t *testing.T) {
	for _, tt := range []struct {
		code int
		text string
	}{
		{Timeout, "Timeout"},
		{NotSupported, "NotSupported"},
		{TemporarilyUnavailable, "TemporarilyUnavailable"},
		{MalformedRequest, "MalformedRequest"},
		{Crash, "Crash"},
		{Abort, "Abort"},
		{KeyDoesNotExist, "KeyDoesNotExist"},
		{KeyAlreadyExists, "KeyAlreadyExists"},
		{PreconditionFailed, "PreconditionFailed"},
		{TxnConflict, "TxnConflict"},
		{1000, "ErrorCode<1000>"},
	} {
		if got, want := ErrorCodeText(tt.code), tt.text; got != want {
			t.Errorf("code %d=%s, want %s", tt.code, got, want)
		}
	}
}

func TestRPCError_ErrorOther(t *testing.T) {
	if got, want := NewRPCError(Crash, "foo").Error(), `RPCError(Crash, "foo")`; got != want {
		t.Fatalf("error=%s, want %s", got, want)
	}
}

func TestRPCError_ErrorCodeOther(t *testing.T) {
	var err error = NewRPCError(Crash, "foo")
	if ErrorCode(err) != Crash {
		t.Fatalf("error=%d, want %d", ErrorCode(err), Crash)
	}

	err = fmt.Errorf("foo: %w", err)
	if ErrorCode(err) != Crash {
		t.Fatalf("error=%d, want %d", ErrorCode(err), Crash)
	}
}