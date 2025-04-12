In this challenge, we're going to create a node. 

In Maelstrom, we create a node which is a binary that receives JSON messages from STDIN and sends JSON messages to STDOUT. The full protocol spec is here:

https://github.com/jepsen-io/maelstrom/blob/main/doc/protocol.md

We're provided with this Maelstrom node library:

https://pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#section-readme

This library provides maelstrom.Node that handles the basic boilerplate for you. It lets you register handler functions for each message type-- similar to how http.Handler works in the standard library.

