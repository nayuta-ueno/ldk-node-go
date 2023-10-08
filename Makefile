all: build

build:
	go build .

run: build
	LD_LIBRARY_PATH=./ldk_node ./ldk-node-go

