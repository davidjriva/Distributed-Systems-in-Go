package main

import (
	"encoding/json"
	"log"
	"errors"
	"sync"
	"strconv"
	"maelstrom-broadcast/safeslice"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

/*
	copyStringMap creates and returns a shallow copy of a map[string]any.
	It assumes the values are safe to copy directly (e.g., strings, numbers, etc.)
	and does not perform a deep copy of nested structures.
*/
func copyStringMap(original map[string]any) map[string]any {
    copyMap := make(map[string]any, len(original))
    for k, v := range original {
        copyMap[k] = v // v is string, so value copy is fine
    }
    return copyMap
}

/*
	extractCurrentNodesNeighbors extracts the list of neighbor node IDs for a given node
	from the "topology" field of the incoming message body.

	Parameters:
		- body: the JSON-decoded message body, expected to contain a "topology" field.
		- nodeId: the ID of the current node.

	Returns:
		- A slice of strings representing the neighbor node IDs.
		- An error if the "topology" field is missing or improperly formatted.
*/
func extractCurrentNodesNeighbors(body map[string]any, nodeId string) ([]string, error) {	
	topologyRaw, ok := body["topology"].(map[string]any)
	if !ok {
		return nil, errors.New("topology field is not a map")
	}

	currNodeNeighbors, ok := topologyRaw[nodeId].([]any)
	if !ok {
		return nil, errors.New("neighbors list not found or invalid")
	}

	neighbors := make([]string, len(currNodeNeighbors))
	for i, v := range currNodeNeighbors {
		neighbors[i], ok = v.(string)

		if !ok {
			return nil, errors.New("neighbor value is not a string")
		}
	}

	return neighbors, nil
}

/*
	broadcastMessageToAllNeighbors sends the given message body to all neighbor nodes.

	Parameters:
	- neighbors: a slice of neighbor node IDs to which the message will be sent.
	- neighborBody: the message payload represented as a map[string]any to send.
	- n: a pointer to the Maelstrom node used to send messages.

	For each neighbor in the slice, the function attempts to send the message.
	If sending fails, the error is logged but the function continues sending to remaining neighbors.
*/
func broadcastMessageToAllNeighbors(neighbors []string, neighborBody map[string]any, n *maelstrom.Node) {
	for _, neighbor := range neighbors {
		if err := n.Send(neighbor, neighborBody); err != nil {
			log.Printf("Error sending to %s: %v", neighbor, err)
		}
	}
}

/*
	This program...
*/
func main() {
    // Initialize a new Maelstrom node for the program to run on.
    n := maelstrom.NewNode()

	// Create a thread-safe slice using the custom 'safeslice' module
	broadcastVals := safeslice.NewSafeSlice()

	var (
		seen sync.Map // A map to track already seen messages with O(1) lookups
		neighbors []string // Define a slice to store the node's neighbors in-memory.
	)

	// Handle the 'broadcast' message type
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err:=json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgFloat, ok := body["message"].(float64)
		if !ok {
			// Handle error: value is not an int
			return errors.New("message field is not a float64")
		}

		/*
			Broadcast message to all neighbors (gossiping)

			1.) Check if we've already seen this message. If we have, then simply reply, otherwise continue.
			2.) Create broadcast message to send to neighboring nodes.
			3.) Send all values to the neighboring node(s).
		*/
		key := strconv.FormatFloat(msgFloat, 'f', -1, 64)

		if _, alreadySeen := seen.Load(key); !alreadySeen {
			// Add the current value in 'message' to our map of already seen values 
			seen.Store(key, true)

			// Safely append the message value to the list of all messages
			broadcastVals.Append(msgFloat)

			// Create a 'broadcast' message to send to all neighbors
			neighborBody := copyStringMap(body)
			neighborBody["type"] = "broadcast"
			neighborBody["message"] = msgFloat

			// Send message to all neighbors
			broadcastMessageToAllNeighbors(neighbors, neighborBody, n)
		}

		// Remove message field from response if it exists
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
		
		// Create a copy of the messages safely
		valsCopy := broadcastVals.GetCopy()

		body["messages"] = valsCopy

		return n.Reply(msg, body)
	})

	// Handle the 'topology' message type
	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err:=json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Extract the neighbors from the "topology" field and store it in memory.
		// The node only stores the neighbors corresponding to this specific node.
		updatedNeighbors, err := extractCurrentNodesNeighbors(body, n.ID())
		if err != nil {
			log.Fatalf("Failed to extract neighbors for node %s: %v", n.ID(), err)
		}

		neighbors = updatedNeighbors

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
