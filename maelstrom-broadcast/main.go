package main

import (
	"encoding/json"
	"log"
	"errors"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func copyStringMap(original map[string]any) map[string]any {
    copyMap := make(map[string]any, len(original))
    for k, v := range original {
        copyMap[k] = v // v is string, so value copy is fine
    }
    return copyMap
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

		// Extract the neighbors from the "topology" field and store it in memory
		topologyRaw, ok := body["topology"].(map[string]any)
		if !ok {
			return errors.New("topology field is not a map")
		}
		
		var topologyNeighbors []string
		for key := range topologyRaw {
			topologyNeighbors = append(topologyNeighbors, key)
		}
		neighbors = topologyNeighbors

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
