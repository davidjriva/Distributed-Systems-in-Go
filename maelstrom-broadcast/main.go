package main

import (
	"encoding/json"
	"log"
	"errors"
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
		neighbors[i], _ = v.(string)
	}

	return neighbors, nil
}

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

	/*
		Define a slice to store the node's neighbors in-memory.
	*/
	neighbors := []string{}

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

		messageAlreadySeen := false

		for _, val := range broadcastVals {
			if val == msgFloat {
				messageAlreadySeen = true
				break
			}
		}

		if !messageAlreadySeen {
			broadcastVals = append(broadcastVals, msgFloat)

			neighborBody := copyStringMap(body)
			neighborBody["type"] = "broadcast"
			neighborBody["message"] = msgFloat

			for _, neighbor := range neighbors {
				if err := n.Send(neighbor, neighborBody); err != nil {
					log.Printf("Error sending to %s: %v", neighbor, err)
				}
			}
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
		body["messages"] = broadcastVals

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
