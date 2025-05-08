# Build latest artifacts
./build.sh

# Run maelstorm tests
num_seconds_to_run=60
node_count=100
/Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/maelstrom test -w unique-ids --bin /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom-unique-ids/maelstrom-unique-ids --time-limit $num_seconds_to_run --rate 1000 --node-count $node_count --availability total --nemesis partition

# Move tests into `/maelstrom` directory to be served
rm -r /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/store

mv ./store /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom