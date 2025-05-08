# Challenge 2: Generate unique IDs

In this challenge, you’ll need to implement a globally-unique ID generation system that runs against Maelstrom’s unique-ids workload. Your service should be totally available, meaning that it can continue to operate even in the face of network partitions.

# Maelstrom Workload: Unique IDs
simple workload for ID generation systems. Clients ask servers to generate an ID, and the server should respond with an ID. The test verifies that those IDs are globally unique.

Your node will receive a request body like:
```json
{"type": "generate",
"msg_id": 2}
```

And should respond with something like:
```json
{"type": "generate_ok",
"in_reply_to": 2,
"id": 123}
```

IDs may be of any type--strings, booleans, integers, floats, compound JSON values, etc.

# RPC: Generate
Asks a node to generate a new ID. Servers respond with a generate_ok message containing an id field, which should be a globally unique value. IDs may be of any type.

Request:
```golang
{:type (eq "generate"), :msg_id Int}
```

Response:
```golang
{:type (eq "generate_ok"),
 :id Any,
 #schema.core.OptionalKey{:k :msg_id} Int,
 :in_reply_to Int}
```