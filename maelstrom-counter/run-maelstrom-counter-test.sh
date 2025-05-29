# Build latest artifacts
./build.sh

# Run maelstorm tests
/Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/maelstrom test -w g-counter --bin /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom-counter/maelstrom-counter --node-count 3 --rate 100 --time-limit 60 --nemesis partition

# Move tests into `/maelstrom` directory to be served
rm -r /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom/store

mv ./store /Users/davidriva/Desktop/Coding-Projects/Distributed-Systems-in-Go/maelstrom