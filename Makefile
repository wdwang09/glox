build:
	go build .

test:
	go test glox/src/test

check-fmt:
	gofmt -l .

all: build test check-fmt
