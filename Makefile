build:
	mkdir -p builds/
	go build -o builds/ cmd/compserv-server.go

test:
	go test -v ./...
