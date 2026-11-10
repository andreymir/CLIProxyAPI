package executor

import (
	"net/http"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/runtime/executor/helps"
)

// enforceClaudeOAuthClaudeCodeOnly is switched off only by an init in this
// package's tests, because upstream tests drive OAuth tokens through other clients.
var enforceClaudeOAuthClaudeCodeOnly = true

// rejectUnconfirmedClaudeOAuth limits Claude subscription tokens to Claude Code.
// It returns 400 because upstream retries 403 and may cool down the credential.
func rejectUnconfirmedClaudeOAuth(apiKey string, detection helps.ClaudeCodeRequestDetection) error {
	if !enforceClaudeOAuthClaudeCodeOnly || !isClaudeOAuthToken(apiKey) {
		return nil
	}
	if detection.Confirmed || (detection.ConfirmedShape && detection.NativeClient) {
		return nil
	}
	return statusErr{code: http.StatusBadRequest, msg: "Claude OAuth credentials accept only Claude Code clients"}
}
