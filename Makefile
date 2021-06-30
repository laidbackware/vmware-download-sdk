all: test

TEST_ARGS ?= -v -count=1

test:
	go test $(TEST_ARGS) ./...

test.debug:
	go test -tags debug $(TEST_ARGS) ./...
