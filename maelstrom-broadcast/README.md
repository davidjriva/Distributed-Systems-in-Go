# Challenge a: Single-Node Broadcast
In this challenge, you’ll need to implement a broadcast system that gossips messages between all nodes in the cluster. Gossiping is a common way to propagate information across a cluster when you don’t need strong consistency guarantees.

This challenge is broken up in multiple sections so that you can build out your system incrementally. First, we’ll start out with a single-node broadcast system. That may sound like an oxymoron but this lets us get our message handlers working correctly in isolation before trying to share messages between nodes.

## Specification
Your node will need to handle the "broadcast" workload which has 3 RPC message types: broadcast, read, & topology. Your node will need to store the set of integer values that it sees from broadcast messages so that they can be returned later via the read message RPC.

The Go library has two methods for sending messages:

1. Send() sends a fire-and-forget message and doesn’t expect a response. As such, it does not attach a message ID.

2. RPC() sends a message and accepts a response handler. The message will be decorated with a message ID so the handler can be invoked when a response message is received.

**Data can be stored in-memory as node processes are not killed by Maelstrom.**

## RPC: `broadcast`

This message requests that a value be broadcast out to all nodes in the cluster. The value is always an integer and it is unique for each message from Maelstrom.

Your node will receive a request message body that looks like this:

```json
{
  "type": "broadcast",
  "message": 1000
}
```

It should store the "message" value locally so it can be read later. In response, it should send an acknowledge with a broadcast_ok message:

```json
{
  "type": "broadcast_ok"
}
```

## RPC: `read`

This message requests that a node return all values that it has seen.

Your node will receive a request message body that looks like this:

```json
{
  "type": "read"
}
```

In response, it should return a read_ok message with a list of values it has seen:

```json
{
  "type": "read_ok",
  "messages": [1, 8, 72, 25]
}
```

**The order of the returned values does not matter.**

## RPC: `topology`

This message informs the node of who its neighboring nodes are. Maelstrom has multiple topologies available or you can ignore this message and make your own topology from the list of nodes in the Node.NodeIDs() method. All nodes can communicate with each other regardless of the topology passed in.

Your node will receive a request message body that looks like this:

```json
{
  "type": "topology",
  "topology": {
    "n1": ["n2", "n3"],
    "n2": ["n1"],
    "n3": ["n1"]
  }
}
```

In response, your node should return a topology_ok message body:
```json
{
  "type": "topology_ok"
}
```

## Solution a:
The solution to this challenge is fairly straightforward. You need to store incoming values in a Go slice and handle messages accordingly. When a read message is received, simply return the current contents of the slice.

# Challenge b: Multi-Node Broadcast

In this challenge, we’ll build on our Single-Node Broadcast implementation and replicate our messages across a cluster that has no network partitions.

## Specification

Your node should propagate values it sees from broadcast messages to the other nodes in the cluster. It can use the topology passed to your node in the topology message or you can build your own topology.

The simplest approach is to simply send a node’s entire data set on every message, however, this is not practical in a real-world system. Instead, try to send data more efficiently as if you were building a real broadcast system.

Values should propagate to all other nodes within a few seconds.

## Solution b

This challenge was more complex, as it required reliably broadcasting messages across the entire network using gossip. When a node receives a new message, it first stores the message locally, then gossips it to all its neighbors. If the message has already been seen, the node simply ignores it and waits for the next one.

# Challenge c: Fault Tolerant Broadcast

In this challenge, we’ll build on our Multi-Node Broadcast implementation, however, this time we’ll introduce network partitions between nodes so they will not be able to communicate for periods of time.