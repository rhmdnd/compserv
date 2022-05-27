build:
	mkdir -p builds/
	go build -o builds/ cmd/compserv-server.go
	go build -o builds/ cmd/migrate/migrate.go
