# Challenge 4: Grow-Only Counter

In this challenge, you’ll need to implement a stateless, grow-only counter which will run against Maelstrom’s g-counter workload. This challenge is different than before in that your nodes will rely on a sequentially-consistent key/value store service provided by Maelstrom.

## Specification
Your node will need to accept two RPC-style message types: add & read. Your service need only be eventually consistent: given a few seconds without writes, it should converge on the correct counter value.

Please note that the final read from each node should return the final & correct count.

## RPC `add`

Your node should accept add requests and increment the value of a single global counter. Your node will receive a request message body that looks like this:

```json
{
  "type": "add",
  "delta": 123
}
```

and it will need to return an "add_ok" acknowledgement message:
```json
{
    "type": "add_ok"
}
```

## RPC `read`

Your node should accept read requests and return the current value of the global counter. Remember that the counter service is only sequentially consistent. Your node will receive a request message body that looks like this:

```json
    "type": "read"
```

and it will need to return a "read_ok" message with the current value:
```json
    "type": "read_ok"
    "value": 1234
```

## Service `seq-kv`

Maelstrom provides a sequentially-consistent key/value store called seq-kv which has read, write, & cas operations. The Go library provides a KV wrapper for this service that you can instantiate with NewSeqKV():

```golang
node := maelstrom.NewNode()
kv := maelstrom.NewSeqKV(node)
```

The API is as follows:

```
func (kv *KV) Read(ctx context.Context, key string) (any, error)
    Read returns the value for a given key in the key/value store. Returns an
    *RPCError error with a KeyDoesNotExist code if the key does not exist.

func (kv *KV) ReadInt(ctx context.Context, key string) (int, error)
    ReadInt reads the value of a key in the key/value store as an int.

func (kv *KV) Write(ctx context.Context, key string, value any) error
    Write overwrites the value for a given key in the key/value store.

func (kv *KV) CompareAndSwap(ctx context.Context, key string, from, to any, createIfNotExists bool) error
    CompareAndSwap updates the value for a key if its current value matches the
    previous value. Creates the key if createIfNotExists is true.

    Returns an *RPCError with a code of PreconditionFailed if the previous value
    does not match. Return a code of KeyDoesNotExist if the key did not exist.
```