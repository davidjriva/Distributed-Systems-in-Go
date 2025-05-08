# Build latest artifacts
./build.sh

# Run maelstorm tests
/Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/maelstrom test -w unique-ids --bin /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom-unique-ids/maelstrom-unique-ids --time-limit 5 --rate 1000 --node-count 3 --availability total --nemesis partition

# Move tests into `/maelstrom` directory to be served
rm -r /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/store

mv ./store /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom