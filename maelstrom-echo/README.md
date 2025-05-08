# Introduction

In this challenge, we're going to create a node. 

The node will receive an "echo" message from Maelstrom that looks like:
```json
{
  "src": "c1",
  "dest": "n1",
  "body": {
    "type": "echo",
    "msg_id": 1,
    "echo": "Please echo 35"
  }
}
```

Nodes and clients are sequentially numbered (i.e. n1,n2,...nk). Nodes are prefixed with "n" and external clients are prefixed with "c". Message IDs are unique per source node which is handled by the Go Maelstrom library.

My job is to send a message with the same body back to the client but with a message type of "echo_ok". It should also associate itself with the original message by setting the "in_reply_to" field to the original message ID. This reply field is handled automatically if you use the `Node.Reply()` method.

The response may look something like:
```json
{
  "src": "n1",
  "dest": "c1",
  "body": {
    "type": "echo_ok",
    "msg_id": 1,
    "in_reply_to": 1,
    "echo": "Please echo 35"
  }
}
```

# Task 1: Implement a node

