SHELL = /usr/bin/env -S bash
.PHONY: b build c clean

b build:
	for file in ./cmd/*; do go build $$file; done

c clean:
	go clean -r -i
	-rm -rvf trem tremd $(wildcard *.test)