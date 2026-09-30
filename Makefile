.PHONY: build-tui build-server build-all run-tui run-server test clean

LDFLAGS := -ldflags="-w -s"

build-tui:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/vault cmd/tui/main.go

build-server:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/vault-server cmd/server/main.go

build-all: build-tui build-server

run-tui:
	go run cmd/tui/main.go

run-server:
	go run cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/
