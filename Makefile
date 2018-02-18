HANDLER_DIRS=$(shell find ./handlers -maxdepth 1 -type d)
BINS=$(patsubst ./handlers/%,./bin/%, $(HANDLER_DIRS))

environment:
	dep ensure
	mkdir -p bin

./bin/%:
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o $@ $(shell find ./handlers/$(shell basename $@) -type f -name '*.go')

clean:
	rm -rf ./bin

build: $(BINS)
