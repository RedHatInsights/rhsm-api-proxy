# rhsm-api-go
A [Caddy](https://github.com/caddyserver/caddy)-based reverse proxy with custom middleware

[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/rhsm-api-proxy)](https://goreportcard.com/report/github.com/RedHatInsights/rhsm-api-proxy) [![Go Reference](https://pkg.go.dev/badge/github.com/RedHatInsights/rhsm-api-proxy.svg)](https://pkg.go.dev/github.com/RedHatInsights/rhsm-api-proxy)

## Modules

### RBAC
The RBAC module is middleware that obtains RBAC access for an authenticated request
and inserts the base64-encoded access list as a request header. 

## Usage
This repo builds a Caddy binary with custom modules included. A Makefile is included
for convenience.

### Test
Execute Go unit tests
```
$ make test
go test ./...
?   	github.com/RedHatInsights/rhsm-api-proxy/cmd/caddy	[no test files]
ok  	github.com/RedHatInsights/rhsm-api-proxy/modules/rbac	(cached)
```

### Run
Run Caddy with the local config example
```
$ make run
go run cmd/caddy/caddy.go run -config config/local.json 
2021/05/04 11:47:50.589	INFO	using provided configuration	{"config_file": "config/local.yml", "config_adapter": "yaml"}
2021/05/04 11:47:50.590	INFO	admin	admin endpoint started	{"address": "tcp/localhost:2019", "enforce_origin": false, "origins": ["localhost:2019", "[::1]:2019", "127.0.0.1:2019"]}
2021/05/04 11:47:50.591	INFO	tls.cache.maintenance	started background certificate maintenance	{"cache": "0xc0002d60e0"}
2021/05/04 11:47:50.591	INFO	tls	cleaned up storage units
2021/05/04 11:47:50.591	INFO	autosaved config	{"file": "/home/peasters/.config/caddy/autosave.json"}
2021/05/04 11:47:50.592	INFO	serving initial configuration
```

### Build
Generate a `caddy` binary with modules embedded
```
$ make build
go build -o bin/caddy cmd/caddy/caddy.go
```

## License
Apache 2.0

See LICENCE to see the full text.
