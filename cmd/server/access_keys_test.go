package main

import (
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
)

func TestHasConfiguredAccessKey(t *testing.T) {
	cases := []struct {
		name string
		keys []string
		want bool
	}{
		{name: "nil", keys: nil, want: false},
		{name: "blank only", keys: []string{"", "   "}, want: false},
		{name: "one real key", keys: []string{"", "key"}, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.APIKeys = tc.keys
			if got := hasConfiguredAccessKey(cfg); got != tc.want {
				t.Fatalf("hasConfiguredAccessKey() = %v, want %v", got, tc.want)
			}
		})
	}
	if hasConfiguredAccessKey(nil) {
		t.Fatal("hasConfiguredAccessKey(nil) = true, want false")
	}
}
