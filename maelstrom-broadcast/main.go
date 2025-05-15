package main

import (
	"encoding/json"
	"log"
	"errors"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

/*
   This program...
*/
func main() {
    // Initialize a new Maelstrom node for the program to run on.
    n := maelstrom.NewNode()

	/* 
		Define a slice to store message values in-memory.
		In golang a slice is similar to an array but can grow and shrink dynamically.
	*/
	broadcastVals := []float64{}

	// Handle the 'broadcast' message type
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err:=json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgInt, ok := body["message"].(float64)
		if !ok {
			// Handle error: value is not an int
			return errors.New("message field is not a float64")
		}

		broadcastVals = append(broadcastVals, msgInt)

		// Remove message field if it exists
		delete(body, "message")

		body["type"] = "broadcast_ok"

		return n.Reply(msg, body)
	})

	// Handle the 'read' message type
	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err:=json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["messages"] = broadcastVals

		return n.Reply(msg, body)
	})

	// Handle the 'topology' message type
	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err:=json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Remove "topology" key if it exists
		delete(body, "topology")

		body["type"] = "topology_ok"

		return n.Reply(msg, body)
	})

    // Start the Maelstrom node, which listens for incoming messages.
    if err := n.Run(); err != nil {
        log.Fatal(err)
    }
}
