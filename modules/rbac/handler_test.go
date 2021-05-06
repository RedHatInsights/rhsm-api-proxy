package rbac

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedHatInsights/rbac-client-go"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockRBAC struct {
	mock.Mock
}

func (m *mockRBAC) GetAccess(ctx context.Context, identity, username string) (rbac.AccessList, error) {
	args := m.Called(ctx, identity, username)
	return args.Get(0).(rbac.AccessList), args.Error(1)
}

func TestServeHTTP(t *testing.T) {
	client := new(mockRBAC)
	m := Middleware{
		rbacClient: client,
	}

	handlerFired := false
	next := caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		handlerFired = true
		assert.NotEmpty(t, r.Header.Get(rbacHeader))
		return nil
	})

	identity := "foo"
	acl := rbac.AccessList{{Permission: "a:b:c"}}
	client.On("GetAccess", mock.Anything, identity, "").Return(acl, nil)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add(identityHeader, identity)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req, next)

	client.AssertExpectations(t)
	assert.True(t, handlerFired, "next handler should fire")
}

func TestServeHTTP_ServiceError(t *testing.T) {
	client := new(mockRBAC)
	m := Middleware{
		rbacClient: client,
		logger:     zap.NewNop(),
	}

	handlerFired := false
	next := caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		handlerFired = true
		return nil
	})

	identity := "foo"
	client.On("GetAccess", mock.Anything, identity, "").Return(rbac.AccessList{}, errors.New(""))

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add(identityHeader, identity)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req, next)

	client.AssertExpectations(t)
	assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode)
	assert.False(t, handlerFired, "next handler should not fire")
}

func TestServeHTTP_NoIdentity(t *testing.T) {
	client := new(mockRBAC)
	m := Middleware{
		rbacClient: client,
		logger:     zap.NewNop(),
	}

	handlerFired := false
	next := caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		handlerFired = true
		return nil
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req, next)

	client.AssertNotCalled(t, "GetAccess")
	assert.True(t, handlerFired, "next handler should fire")
}
