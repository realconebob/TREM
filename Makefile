.PHONY: c clean r run t test

r run:
	go run

t test:
	go test

c clean:
	go clean -r -i
	-rm -rv *.test