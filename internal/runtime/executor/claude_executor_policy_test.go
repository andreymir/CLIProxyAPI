package executor

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/runtime/executor/helps"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v8/sdk/translator"
)

func init() {
	enforceClaudeOAuthClaudeCodeOnly = false
}

func enableClaudeOAuthGate(t *testing.T) {
	t.Helper()
	enforceClaudeOAuthClaudeCodeOnly = true
	t.Cleanup(func() { enforceClaudeOAuthClaudeCodeOnly = false })
}

func requireBadRequest(t *testing.T, err error) {
	t.Helper()
	var status statusErr
	if !errors.As(err, &status) {
		t.Fatalf("err = %v, want statusErr", err)
	}
	if status.code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", status.code, http.StatusBadRequest)
	}
}

func TestClaudeOAuthRejectsUnconfirmedClient(t *testing.T) {
	enableClaudeOAuthGate(t)
	const oauth = "sk-ant-oat-test"

	requireBadRequest(t, rejectUnconfirmedClaudeOAuth(oauth, helps.ClaudeCodeRequestDetection{}))
	requireBadRequest(t, rejectUnconfirmedClaudeOAuth(oauth, helps.ClaudeCodeRequestDetection{ConfirmedShape: true}))

	allowed := []helps.ClaudeCodeRequestDetection{
		{Confirmed: true},
		{ConfirmedShape: true, NativeClient: true},
	}
	for _, detection := range allowed {
		if err := rejectUnconfirmedClaudeOAuth(oauth, detection); err != nil {
			t.Fatalf("detection %+v: err = %v, want nil", detection, err)
		}
	}
	if err := rejectUnconfirmedClaudeOAuth("sk-ant-api03-test", helps.ClaudeCodeRequestDetection{}); err != nil {
		t.Fatalf("API key: err = %v, want nil", err)
	}
}

func claudeGateFixture(t *testing.T) (*ClaudeExecutor, *cliproxyauth.Auth, cliproxyexecutor.Request, cliproxyexecutor.Options, *atomic.Int32) {
	t.Helper()
	enableClaudeOAuthGate(t)
	var hits atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(upstream.Close)
	auth := &cliproxyauth.Auth{
		ID:         t.Name(),
		Provider:   "claude",
		Attributes: map[string]string{"api_key": "sk-ant-oat-test", "base_url": upstream.URL},
		Metadata:   map[string]any{"account_uuid": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"},
	}
	request := cliproxyexecutor.Request{
		Model:   "claude-sonnet-4-5",
		Payload: []byte(`{"model":"claude-sonnet-4-5","max_tokens":1,"messages":[{"role":"user","content":"hi"}]}`),
	}
	options := cliproxyexecutor.Options{SourceFormat: sdktranslator.FromString("claude")}
	return NewClaudeExecutor(&config.Config{}), auth, request, options, &hits
}

func requireUpstreamNotCalled(t *testing.T, hits *atomic.Int32) {
	t.Helper()
	if got := hits.Load(); got != 0 {
		t.Fatalf("upstream hits = %d, want 0", got)
	}
}

func TestClaudeExecuteRejectsUnconfirmedOAuth(t *testing.T) {
	executor, auth, request, options, hits := claudeGateFixture(t)
	_, err := executor.Execute(context.Background(), auth, request, options)
	requireBadRequest(t, err)
	requireUpstreamNotCalled(t, hits)
}

func TestClaudeExecuteStreamRejectsUnconfirmedOAuth(t *testing.T) {
	executor, auth, request, options, hits := claudeGateFixture(t)
	options.Stream = true
	_, err := executor.ExecuteStream(context.Background(), auth, request, options)
	requireBadRequest(t, err)
	requireUpstreamNotCalled(t, hits)
}

func TestClaudeCountTokensUpstreamRejectsUnconfirmedOAuth(t *testing.T) {
	executor, auth, request, options, hits := claudeGateFixture(t)
	_, err := executor.countTokensUpstream(context.Background(), auth, request, options)
	requireBadRequest(t, err)
	requireUpstreamNotCalled(t, hits)
}
