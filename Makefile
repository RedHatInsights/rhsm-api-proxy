test:
	go test ./...

build:
	go build -o bin/caddy cmd/caddy/caddy.go

run:
	go run cmd/caddy/caddy.go run --config config/local.json
