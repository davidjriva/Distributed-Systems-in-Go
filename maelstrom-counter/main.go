package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

/*
This program...
*/
func main() {
	// Initialize a new Maelstrom node for the program to run on.
	n := maelstrom.NewNode()

	// Initialize a key-value store to persist operations on even in the case of node failures.
	kv := maelstrom.NewSeqKV(n)

	// Handle the 'add' message type
	n.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Read in request value 'delta' as an int
		deltaFloat, ok := body["delta"].(float64)
		if !ok {
			return fmt.Errorf("delta is not a number")
		}
		delta := int(deltaFloat)

		/*
			Make a write to the key belonging to this node.
			Utilize the node ID so that nodes have less competition for writes.
		*/
		key := fmt.Sprintf("counter-%s", n.ID())
		for {
			// Read the current value of the global counter "g_ct" from the kv store
			curr_ct, err := kv.ReadInt(context.TODO(), key)
			if rpcErr, ok := err.(*maelstrom.RPCError); ok && rpcErr.Code == maelstrom.KeyDoesNotExist {
				curr_ct = 0
			} else if err != nil {
				return err
			}

			new_val := curr_ct + delta

			err = kv.CompareAndSwap(context.TODO(), key, curr_ct, new_val, true)

			if err == nil {
				// Write succeeded
				break
			}

			if rpcErr, ok := err.(*maelstrom.RPCError); !ok || rpcErr.Code != maelstrom.PreconditionFailed {
				return err // Unrecoverable error
			}

			// Write failed, retry
		}

		// Remove message field from response if it exists
		res := map[string]any{
			"type":        "add_ok",
			"msg_id":      body["msg_id"],
			"in_reply_to": body["in_reply_to"],
		}

		return n.Reply(msg, res)
	})

	// Handle the 'read' message type
	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		total := 0

		for _, id := range n.NodeIDs() {
			key := fmt.Sprintf("counter-%s", id)

			val, err := kv.ReadInt(context.TODO(), key)
			if err != nil {
				return err
			}

			total += val
		}

		// Remove message field from response if it exists
		res := map[string]any{
			"type":        "read_ok",
			"value":       total,
			"msg_id":      body["msg_id"],
			"in_reply_to": body["in_reply_to"],
		}

		return n.Reply(msg, res)
	})

	// Start the Maelstrom node, which listens for incoming messages.
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
