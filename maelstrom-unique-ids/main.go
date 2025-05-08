package main

import (
    "encoding/json"
	"strconv"
    "log"
	"math"
	"sync/atomic"
	"time"
	"math/rand"
	"fmt"
    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	// Create a Maelstrom node
	n := maelstrom.NewNode()

	// Create an atomic counter (integer) to track unique IDs for this node
	var GUID atomic.Uint64

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update the message type
		body["type"] = "generate_ok"

		// Atomically increment and get the new unique ID
		id := GUID.Add(1) - 1 // `Add(1)` returns the *new* value, so subtract 1 for the *previous* value

		// Handle cases where id reaches its maximum possible value and needs to be reset to zero
		if id == math.MaxUint64 {
			log.Println("GUID has reached the maximum value!")
			// Reset GUID back to zero
			GUID.Store(0)
		}
	
		// Set the ID as a unique key across all Node instances.
		// For node n1, the ID could be n1_0
		// For node n2, the ID could be n2_0
		uniqueID := n.ID() + "_" + strconv.FormatUint(id, 10)

		// Append the current timestamp to maintain uniqueness as long as the process runs
		currentTime := uint64(time.Now().UnixNano())
		uniqueID += "_" + strconv.FormatUint(currentTime, 10)

		// Add random string to ensure uniqueness even in the case of the currentTime variable overflowing
		randomBytes := make([]byte, 8)
		_, err := rand.Read(randomBytes)

		if err != nil {
			log.Fatal("Error generating random bytes:", err)
		}

		randomString := fmt.Sprintf("%x", randomBytes)

		uniqueID += "_" + randomString

		// Store the uniqueID on the response's ID field.
		body["id"] = uniqueID

		// Send the response
		return n.Reply(msg, body)
	})	

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}