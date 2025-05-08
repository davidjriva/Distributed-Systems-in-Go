package main

import (
    "encoding/json"
    "log"
	"sync"
	"sync/atomic"
    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	// Create a Maelstrom node
	n := maelstrom.NewNode()

	// Create an atomic counter (integer)
	var uid atomic.Uint64

	// Create a wait group
	// Wait groups help us wait for all go routines to finish their work
	var wg sync.WaitGroup

	// Define the "generate" response logic
	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as a loosely-typed map
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back
		body["type"] = "generate_ok"
		
		// Synchronize on atomic integer
		wg.Add(1)

		// Start new goroutine to update the "id" variable
		// and increase the count of the uid variable by one.
		go func(){
			// Assign a unique ID
			body["id"] = uid.Load()

			// Increase the count by one
			uid.Add(1)

			// Free lock
			wg.Done()
		}()

		// Wait for operation to finish
		wg.Wait()

		// Echo the original message back with the updated message type
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
