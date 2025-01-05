
run:
	go run cmd/d9s/main.go
	
clean:
	rm -rf bin

build:
	go build -o ./bin/d9s cmd/d9s/main.go
