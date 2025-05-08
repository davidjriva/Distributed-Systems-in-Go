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

func generateRandomBytes() []byte {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)

	if err != nil {
		log.Fatal("Error generating random bytes:", err)
	}
	
	return randomBytes
}

func createGUID(nodeID string, id uint64) string {
	// Set the ID as a unique key across all Node instances.
	// For node n1, the ID could be n1_0
	// For node n2, the ID could be n2_0
	uniqueID := nodeID + "_" + strconv.FormatUint(id, 10)

	// Append the current timestamp to maintain uniqueness as long as the process runs
	currentTime := uint64(time.Now().UnixNano())
	uniqueID += "_" + strconv.FormatUint(currentTime, 10)

	// Add random string to ensure uniqueness even in the case of the currentTime variable overflowing
	randomBytes := generateRandomBytes()

	randomString := fmt.Sprintf("%x", randomBytes)

	uniqueID += "_" + randomString

	return uniqueID
}

func incrementGUIDCount(GUID atomic.Uint64) uint64 {
	previousValue := GUID.Add(1) - 1

	if previousValue == math.MaxUint64 - 1 {
		log.Println("GUID has reached the maximum value!")
		// Reset GUID back to zero
		GUID.Store(0)
	}

	return previousValue
}

/*
    This program generates Globally Unique Identifiers (GUIDs) by combining the node's ID, 
    an atomic counter, the current timestamp in nanoseconds, and 64 randomly generated bits. 
    These components together ensure the creation of globally unique IDs for each node.
*/
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

		// Atomically increment GUID variable and get its previous value
		id := incrementGUIDCount(GUID)

		// Programatically create a GUID using the Node's ID, the atomic counter GUID, the timestamp, and 64 random bits
		uniqueID := createGUID(n.ID(), id)

		// Store the uniqueID on the response's ID field.
		body["id"] = uniqueID

		// Send the response
		return n.Reply(msg, body)
	})	

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}