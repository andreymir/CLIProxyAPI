package main

import (
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
)

// hasConfiguredAccessKey reports whether at least one non-blank client key exists.
func hasConfiguredAccessKey(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	for _, key := range cfg.APIKeys {
		if strings.TrimSpace(key) != "" {
			return true
		}
	}
	return false
}
