SHELL = /usr/bin/env -S bash
.PHONY: c clean r run t test b build

b build:
	for file in ./cmd/*; do go build $$file; done

c clean:
	go clean -r -i
	-rm -rvf trem tremd $(wildcard *.test)