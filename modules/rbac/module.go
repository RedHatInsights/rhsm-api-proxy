package rbac

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RedHatInsights/rbac-client-go"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Middleware{})
}

type authorizationService interface {
	GetAccess(ctx context.Context, identity string, username string) (rbac.AccessList, error)
}

// Middleware implements an HTTP handler that requests
// RBAC permissions and inserts them as a header
type Middleware struct {
	// Base URL for the RBAC service
	ServiceURL string `json:"service_url"`
	// Timeout for RBAC service requests
	Timeout caddy.Duration `json:"timeout"`
	// Application name is used when enumerating permissions for a principal
	Application string `json:"application"`

	rbacClient authorizationService
	logger     *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.rbac",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	if len(m.ServiceURL) == 0 {
		return fmt.Errorf("service URL must be specified")
	}
	if len(m.Application) == 0 {
		return fmt.Errorf("application name must be specified")
	}

	m.rbacClient = &rbac.Client{
		HTTPClient:  &http.Client{Timeout: time.Duration(m.Timeout)},
		BaseURL:     m.ServiceURL,
		Application: m.Application,
	}

	m.logger = ctx.Logger(m)

	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	if m.rbacClient == nil {
		return fmt.Errorf("no rbac client")
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ authorizationService        = (*rbac.Client)(nil)
)
