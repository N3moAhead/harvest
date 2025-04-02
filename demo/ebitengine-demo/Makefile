.DEFAULT_GOAL = dev

fmt:
	go fmt cmd/test/main.go

vet: fmt
	go vet cmd/test/main.go

dev: vet
	go run cmd/test/main.go

build: vet
	go build -o ebitengine-test cmd/test/main.go

start: build
	./ebitengine-test

clean:
	rm -f ./ebitengine-test
