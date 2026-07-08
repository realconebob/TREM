SHELL = /usr/bin/env -S bash
.PHONY: c clean r run t test b build

b build:
	for file in ./cmd/*; do go build $$file; done

r run:
	go run

t test:
	go test

c clean:
	go clean -r -i
	-rm -rv *.test