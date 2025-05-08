package main

import (
    "encoding/json"
	"strconv"
    "log"
	"sync/atomic"
    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	// Create a Maelstrom node
	n := maelstrom.NewNode()

	// Create an atomic counter (integer)
	var uid atomic.Uint64

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update the message type
		body["type"] = "generate_ok"
	
		// Atomically increment and get the new unique ID
		id := uid.Add(1) - 1 // `Add(1)` returns the *new* value, so subtract 1 for the *previous* value
	
		// Set the ID
		body["id"] = n.ID() + "_" + strconv.FormatUint(id, 10)
	
		// Send the response
		return n.Reply(msg, body)
	})	

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}