package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Caddy modules
	_ "github.com/RedHatInsights/rhsm-api-proxy/modules/clowder"
	_ "github.com/RedHatInsights/rhsm-api-proxy/modules/rbac"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

func main() {
	// TODO: possibly restructure this and merge clowder settings with config
	caddycmd.Main()
}
