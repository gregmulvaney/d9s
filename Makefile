run:
	go run cmd/d9s/main.go

build: 
	go build -o bin/d9s cmd/d9s/main.go

test: 
	go test ./...
