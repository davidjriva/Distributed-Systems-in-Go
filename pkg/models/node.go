package maelstrom

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

/*
	Node represents a single node in the network.
*/
type Node struct {
	mu sync.Mutex
	wg sync.WaitGroup

	id        string
	nodeIDs   []string
	nextMsgID int

	handlers  map[string]HandlerFunc
	callbacks map[int]HandlerFunc

	// Stdin is for reading messages in from the Maelstrom network.
	Stdin io.Reader

	// Stdout is for writing messages out to the Maelstrom network.
	Stdout io.Writer
}

/*
	NewNode returns a new instance of Node connected to STDIN/STDOUT.
*/
func NewNode() *Node {
	return &Node{
		handlers:  make(map[string]HandlerFunc),
		callbacks: make(map[int]HandlerFunc),

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
}

/*
	Init is used for initializing the node. This is normally called after
	receiving an "init" message but it can also be called manually when
	initializing unit tests.
*/
func (n *Node) Init(id string, nodeIDs []string) {
	n.id = id
	n.nodeIDs = nodeIDs
}

/*
	ID returns the identifier for this node.
	Only valid after "init" message has been received.
*/
func (n *Node) ID() string {
	return n.id
}

/*
	NodeIDs returns a list of all node IDs in the cluster. This list include the
	local node ID and is the same order across all nodes. Only valid after "init"
	message has been received.
*/
func (n *Node) NodeIDs() []string {
	return n.nodeIDs
}

/*
	Handle registers a message handler for a given message type. Will panic if
	registering multiple handlers for the same message type.'
*/
func (n *Node) Handle(typ string, fn HandlerFunc) {
	if _, ok := n.handlers[typ]; ok {
		panic(fmt.Sprintf("duplicate message handler for %q message type", typ))
	}
	n.handlers[typ] = fn
}

/*
	Run executes the main event handling loop. It reads in messages from STDIN
	and delegates them to the appropriate registered handler. This should be
	the last function executed by main().
*/
func (n *Node) Run() error {
	scanner := bufio.NewScanner(n.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()

		// Parse next line from STDIN as a JSON-formatted message.
		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			return fmt.Errorf("unmarshal message: %w", err)
		}

		var body MessageBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return fmt.Errorf("unmarshal message body: %w", err)
		}
		log.Printf("Received %s", msg)

		// What handler should we use for this message?
		if body.InReplyTo != 0 {
			// Extract callback, if replying to a previous message.
			n.mu.Lock()
			h := n.callbacks[body.InReplyTo]
			delete(n.callbacks, body.InReplyTo)
			n.mu.Unlock()

			// If no callback exists, just log a message and skip.
			if h == nil {
				log.Printf("Ignoring reply to %d with no callback", body.InReplyTo)
				continue
			}

			// Handle callback in a separate goroutine.
			n.wg.Add(1)
			go func() {
				defer n.wg.Done()
				n.handleCallback(h, msg)
			}()
			continue
		}

		// If this is not a callback, ensure that a handler is registered.
		var h HandlerFunc
		if body.Type == "init" {
			h = n.handleInitMessage // wraps init message with special handling.
		} else if h = n.handlers[body.Type]; h == nil {
			return fmt.Errorf("No handler for %s", line)
		}

		// Handle message in a separate goroutine.
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			n.handleMessage(h, msg)
		}()
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Wait for all in-flight handlers to complete.
	n.wg.Wait()

	return nil
}

/*
	handleCallback sends msg response to a callback function. Logs error, if one occurs.
*/
func (n *Node) handleCallback(h HandlerFunc, msg Message) {
	if err := h(msg); err != nil {
		log.Printf("callback error: %s", err)
	}
}

/*
	handleMessage sends msg to a handler function. Sends an RPC error if an error is returned.
*/
func (n *Node) handleMessage(h HandlerFunc, msg Message) {
	if err := h(msg); err != nil {
		switch err := err.(type) {
		case *RPCError:
			if err := n.Reply(msg, err); err != nil {
				log.Printf("reply error: %s", err)
			}
		default:
			log.Printf("Exception handling %#v:\n%s", msg, err)
			if err := n.Reply(msg, NewRPCError(Crash, err.Error())); err != nil {
				log.Printf("reply error: %s", err)
			}
		}
	}
}

/*
	Handles init message by initializing a node with ID and nodeIDs in the network
*/
func (n *Node) handleInitMessage(msg Message) error {
	var body InitMessageBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("unmarshal init message body: %w", err)
	}
	n.Init(body.NodeID, body.NodeIDs)

	// Delegate to application initialization handler, if specified.
	if h := n.handlers["init"]; h != nil {
		if err := h(msg); err != nil {
			return err
		}
	}

	// Send back a response that the node has been initialized.
	log.Printf("Node %s initialized", n.id)
	return n.Reply(msg, MessageBody{Type: "init_ok"})
}

/*
	Reply replies to a request with a response body.
*/
func (n *Node) Reply(req Message, body any) error {
	// Extract the message ID from the original message.
	var reqBody MessageBody
	if err := json.Unmarshal(req.Body, &reqBody); err != nil {
		return err
	}

	// We have to marshal/unmarshal to inject our reply message ID.
	b := make(map[string]any)
	if buf, err := json.Marshal(body); err != nil {
		return err
	} else if err := json.Unmarshal(buf, &b); err != nil {
		return err
	}
	b["in_reply_to"] = reqBody.MsgID

	return n.Send(req.Src, b)
}

/*
	Send sends a message body to a given destination node.
*/
func (n *Node) Send(dest string, body any) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(Message{
		Src:  n.id,
		Dest: dest,
		Body: bodyJSON,
	})
	if err != nil {
		return err
	}

	// Synchronize access to STDOUT.
	n.mu.Lock()
	defer n.mu.Unlock()

	log.Printf("Sent %s", buf)

	if _, err = n.Stdout.Write(buf); err != nil {
		return err
	}
	_, err = n.Stdout.Write([]byte{'\n'})
	return err
}

/*
	RPC sends an async RPC request. Handler invoked when response message received.
*/
func (n *Node) RPC(dest string, body any, handler HandlerFunc) error {
	n.mu.Lock()

	// Generate a unique message ID.
	n.nextMsgID++
	msgID := n.nextMsgID

	// Register a handler for our callback.
	n.callbacks[msgID] = handler

	n.mu.Unlock()

	// We have to marshal/unmarshal to inject our message ID.
	b := make(map[string]any)
	if buf, err := json.Marshal(body); err != nil {
		return err
	} else if err := json.Unmarshal(buf, &b); err != nil {
		return err
	}
	b["msg_id"] = msgID

	return n.Send(dest, b)
}

/*
	SyncRPC sends a synchronous RPC request. Returns the response message. RPC
	errors in the message body are converted to *RPCError and are returned.
*/
func (n *Node) SyncRPC(ctx context.Context, dest string, body any) (Message, error) {
	respCh := make(chan Message)
	if err := n.RPC(dest, body, func(m Message) error {
		respCh <- m
		return nil
	}); err != nil {
		return Message{}, err
	}

	// Wait for either the context to finish or for the response message to arrive.
	select {
	case <-ctx.Done():
		return Message{}, ctx.Err()

	case m := <-respCh:
		if err := m.RPCError(); err != nil {
			return m, err
		}
		return m, nil
	}
}

// HandlerFunc is the function signature for a message handler.
type HandlerFunc func(msg Message) error