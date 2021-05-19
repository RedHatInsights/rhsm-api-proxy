package rbac

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

const identityHeader = "X-RH-Identity"
const rbacHeader = "X-RH-RBAC"

// ServeHTTP implements caddyhttp.MiddlewareHandler
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Fetch RBAC if identity header is present
	identity := r.Header.Get(identityHeader)
	if len(identity) > 0 {
		acl, err := m.rbacClient.GetAccess(r.Context(), identity, "")
		// Respond with 503 if RBAC is unavailable
		if err != nil {
			m.logger.Sugar().Errorw("failed to get access for request", "error", err)
			http.Error(w, "Authorization service unavailable", http.StatusServiceUnavailable)
			return nil
		}

		// Insert RBAC header
		aclJSON, err := json.Marshal(acl)
		if err != nil {
			m.logger.Sugar().Errorw("failed to marshal access list to JSON", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return nil
		}
		r.Header.Add(rbacHeader, base64.StdEncoding.EncodeToString(aclJSON))
	}
	return next.ServeHTTP(w, r)
}
