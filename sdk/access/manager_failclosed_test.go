package access

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticateFailsClosedWithoutProviders(t *testing.T) {
	manager := NewManager()
	request := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)

	result, err := manager.Authenticate(context.Background(), request)

	if result != nil {
		t.Fatalf("result = %+v, want nil", result)
	}
	if err == nil {
		t.Fatal("err = nil, want a no-credentials error")
	}
	if err.Code != AuthErrorCodeNoCredentials {
		t.Fatalf("code = %q, want %q", err.Code, AuthErrorCodeNoCredentials)
	}
	if err.HTTPStatusCode() != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", err.HTTPStatusCode(), http.StatusUnauthorized)
	}
}
