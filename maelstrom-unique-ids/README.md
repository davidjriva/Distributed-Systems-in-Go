# Table of Contents

1. [Challenge: Generate Unique IDs](#challenge-2-generate-unique-ids)
    - Overview of the Challenge
    - Service Requirements
    - Description of the Workload
    - Request and Response Example
    - Request Format
    - Response Format

2. [Scripts](#scripts)
    - Useful scripts

3. [Motivation](#motivation)
   - Why GUIDs Are Useful
   - Benefits in Distributed Systems

4. [Solution](#solution)
   - Initial Approach and Issues
   - Refining the Approach for Global Uniqueness
   - Improvements to Handle GUID Overflow
     1. GUID Reset
     2. Timestamp Integration
   - Strengthening with Random Bytes for Future Robustness


# Challenge: Generate unique IDs

In this challenge, you’ll need to implement a globally-unique ID generation system that runs against Maelstrom’s unique-ids workload. Your service should be totally available, meaning that it can continue to operate even in the face of network partitions.

## Maelstrom Workload: Unique IDs
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

## RPC: Generate
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

# Scripts
I created two scripts to help in automating the build process and testing of `main.go`.

* **build.sh**: A utility script that compiles the main.go application, using the -o flag to specify the output binary.

* **run-maelstrom-unique-ids-test.sh**: This script builds the application, runs all the available Maelstrom tests, and saves the /store test outputs into the Maelstrom package for later serving.

# Motivation
GUIDs (Globally Unique Identifiers) are useful because they provide a way to generate unique IDs across distributed systems, ensuring that each ID is distinct even across multiple nodes or servers. This is crucial for maintaining data integrity and preventing collisions in large-scale systems, where multiple entities may need to generate or reference identifiers independently. GUIDs help avoid conflicts without the need for centralized coordination, making them ideal for decentralized or distributed applications.

# Solution
This problem may seem straightforward at first—generate unique IDs across all nodes (servers) running `main.go`.

My initial approach was to have each `Node` return an atomic integer as a unique ID. However, this led to issues as all nodes ended up returning the same ID. While the `GUID` variable was accessed in a thread-safe manner, ensuring unique IDs per node, the IDs were not globally unique across all nodes.

To address this, I refined my approach by leveraging the `Node.id` field to generate globally unique IDs. I concatenated the node's `GUID` with its `Node.id` to create distinct IDs. For example, a `Node` with ID N1 would generate a GUID like N1_0, while another node, N2, would generate N2_0. This approach successfully passed the available tests and generated unique IDs for each node, but it wasn’t foolproof.

The potential issue was that once the GUID reached its maximum value and overflowed, the system would crash, making the server unavailable.

To resolve this, I implemented two key improvements:
1. **GUID Reset**: When the `GUID` reaches its maximum value, it resets to zero, thus preventing any overflow issues.
2. **Timestamp Integration**: I added the current time in nanoseconds (represented as a uint64) to the `GUID`, further ensuring uniqueness. This timestamp ensures that GUIDs can continue to be generated without collision for the next 584 years, until the year 2554.

To further strengthen the system's robustness, I included a third component—64 randomly generated bits. This `randomBytes` term will help prevent collisions once the `currentTime` variable overflows in the year 2554, ensuring that unique IDs can still be generated even in the distant future.

These changes together provide a reliable and future-proof solution for generating globally unique IDs across all nodes.