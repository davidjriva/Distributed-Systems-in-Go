package main

import (
    "encoding/json"
    "strconv"
    "log"
    "math"
    "sync/atomic"
    "time"
    "crypto/rand"
    "fmt"
    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// generateRandomBytes generates 8 random bytes to ensure uniqueness for the GUID.
func generateRandomBytes() []byte {
    randomBytes := make([]byte, 8)
    _, err := rand.Read(randomBytes)

    if err != nil {
        log.Fatal("Error generating random bytes:", err)
    }

    return randomBytes
}

// createGUID creates a Globally Unique Identifier (GUID) based on the node's ID, 
// an atomic counter, the current timestamp, and 64 random bits to ensure uniqueness.
func createGUID(nodeID string, id uint64) string {
    // Combine nodeID and ID to create a base unique key. 
    // Example: For node n1, the ID could be n1_0; for node n2, it could be n2_0.
    uniqueID := nodeID + "_" + strconv.FormatUint(id, 10)

    // Append the current timestamp in nanoseconds to further ensure uniqueness during runtime.
    currentTime := uint64(time.Now().UnixNano())
    uniqueID += "_" + strconv.FormatUint(currentTime, 10)

    // Generate random bytes and convert them to a hexadecimal string to prevent collisions, 
    // even in cases where the timestamp might overflow.
    randomBytes := generateRandomBytes()
    randomString := fmt.Sprintf("%x", randomBytes)

    uniqueID += "_" + randomString

    return uniqueID
}

// incrementGUIDCount atomically increments the GUID counter and handles overflow.
func incrementGUIDCount(GUID *atomic.Uint64) uint64 {
    // Atomically increment the counter and return the previous value.
    previousValue := GUID.Add(1) - 1

    // If the GUID counter reaches the maximum value, reset it to zero to prevent overflow.
    if previousValue == math.MaxUint64-1 {
        log.Println("GUID has reached the maximum value!")
        GUID.Store(0)
    }

    return previousValue
}

/*
    This program generates Globally Unique Identifiers (GUIDs) by combining the node's ID, 
    an atomic counter, the current timestamp in nanoseconds, and 64 randomly generated bits. 
    These components together ensure the creation of globally unique IDs for each node, even across distributed systems.
*/
func main() {
    // Initialize a new Maelstrom node for the program to run on.
    n := maelstrom.NewNode()

    // Create an atomic counter to track unique GUIDs for this particular node.
    var GUID atomic.Uint64

    // Handle the "generate" message type by responding with a new unique GUID.
    n.Handle("generate", func(msg maelstrom.Message) error {
        var body map[string]any
        if err := json.Unmarshal(msg.Body, &body); err != nil {
            return err
        }
    
        // Update the response type to indicate successful GUID generation.
        body["type"] = "generate_ok"

        // Atomically increment the GUID counter and get the previous value.
        id := incrementGUIDCount(&GUID)

        // Create a unique GUID by combining the node ID, the atomic counter, the timestamp, and random bytes.
        uniqueID := createGUID(n.ID(), id)

        // Add the generated uniqueID to the response body.
        body["id"] = uniqueID

        // Send the response back to the requester.
        return n.Reply(msg, body)
    })    

    // Start the Maelstrom node, which listens for incoming messages.
    if err := n.Run(); err != nil {
        log.Fatal(err)
    }
}
