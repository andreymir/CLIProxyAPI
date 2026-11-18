package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	sdkaccess "github.com/router-for-me/CLIProxyAPI/v8/sdk/access"
	"github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

// The fork does not expose the upstream /v1/tunnel WebSocket relay.
func TestServerWebsocketTunnelRouteRemoved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		SDKConfig: config.SDKConfig{
			APIKeys: []string{"tunnel-key"},
		},
	}
	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	apiSrv := NewServer(cfg, auth.NewManager(nil, nil, nil), sdkaccess.NewManager(), cfgPath, WithRequestLoggerFactory(nil))
	t.Cleanup(func() {
		_ = apiSrv.Stop(context.Background())
	})

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/v1/tunnel?role=client&uuid=u1&name=n1", nil)
		req.Header.Set("Authorization", "Bearer tunnel-key")
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
		rec := httptest.NewRecorder()
		apiSrv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s /v1/tunnel = %d, want %d", method, rec.Code, http.StatusNotFound)
		}
	}
}
