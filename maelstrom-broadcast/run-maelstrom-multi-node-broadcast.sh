# Build latest artifacts
./build.sh

# Run maelstorm tests
/Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/maelstrom test -w broadcast --bin /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom-broadcast/maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10

# Move tests into `/maelstrom` directory to be served
rm -r /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/store

mv ./store /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom