.PHONY: bin test all fmt deploy docs server cli setup

all: fmt bin

fmt:
	-go fmt ./...

bin: server cli

server:

cli:
	(cd ./cmd/mcserv; go build)

run: cli
	(cd ./cmd/mcserv; ov site)

devrun:
	-reflex -r '\.go$\' -s make run

devtest:
	-reflex -r '\.go$\' -s make test 

dep:
	dep ensure

test:
	-go test -v ./...

docs:
	./makedocs.sh

setup:
	-go get -u github.com/cespare/reflex
	-dep ensure
