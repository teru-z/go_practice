.DEFAULT_GOAL := run

.PHONY: fmt vet clean build run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

clean: vet
	rm -rf bin/

build: clean
	@mkdir -p bin
	go build -o bin/hello .

run: build
	./bin/hello