package plugin

import (
	"os"
	"testing"
)

func TestIsPathBlocked(t *testing.T) {
	tests := []struct {
		name        string
		blockList   string
		path        string
		shouldBlock bool
	}{
		{
			name:        "single term matches",
			blockList:   "secret",
			path:        "/some/path/secret-database.db",
			shouldBlock: true,
		},
		{
			name:        "multiple terms, middle matches",
			blockList:   "secret,admin,sensitive",
			path:        "/some/root/path/admin/test-check-db",
			shouldBlock: true,
		},
		{
			name:        "block list with spaces around terms",
			blockList:   " secret , config , admin ",
			path:        "/some/root/path/config-file.db",
			shouldBlock: true,
		},
		{
			name:        "no terms match",
			blockList:   "secret,admin,sensitive",
			path:        "/some/path/public-data.db",
			shouldBlock: false,
		},
		{
			name:        "empty block list",
			blockList:   "",
			path:        "/some/path/secret-database.db",
			shouldBlock: false,
		},
		{
			name:        "block list with only commas",
			blockList:   ",,,",
			path:        "/some/path/secret-database.db",
			shouldBlock: false,
		},
		{
			name:        "case insensitive matching",
			blockList:   "Secret",
			path:        "/some/path/secret-database.db",
			shouldBlock: true,
		},
		{
			name:        "partial path matching",
			blockList:   "tmp",
			path:        "/tmp/database.db",
			shouldBlock: true,
		},
		// Default blocklist tests (no env var needed)
		{
			name:        "default blocklist blocks .aws",
			blockList:   "",
			path:        "/home/user/.aws/credentials.db",
			shouldBlock: true,
		},
		{
			name:        "default blocklist blocks .ssh",
			blockList:   "",
			path:        "/home/user/.ssh/keys.db",
			shouldBlock: true,
		},
		{
			name:        "default blocklist blocks grafana.db",
			blockList:   "",
			path:        "/var/lib/grafana/grafana.db",
			shouldBlock: true,
		},
		{
			name:        "default blocklist is case insensitive",
			blockList:   "",
			path:        "/home/user/.AWS/credentials.db",
			shouldBlock: true,
		},
		{
			name:        "default blocklist blocks /etc/shadow",
			blockList:   "",
			path:        "/etc/shadow",
			shouldBlock: true,
		},
		{
			name:        "default blocklist blocks id_rsa",
			blockList:   "",
			path:        "/home/user/.ssh/id_rsa",
			shouldBlock: true,
		},
		{
			name:        "safe path not blocked by default blocklist",
			blockList:   "",
			path:        "/var/data/myapp/analytics.db",
			shouldBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			originalValue := os.Getenv("GF_PLUGIN_BLOCK_LIST")
			defer func() {
				if originalValue == "" {
					_ = os.Unsetenv("GF_PLUGIN_BLOCK_LIST")
				} else {
					_ = os.Setenv("GF_PLUGIN_BLOCK_LIST", originalValue)
				}
			}()

			if tt.blockList == "" {
				_ = os.Unsetenv("GF_PLUGIN_BLOCK_LIST")
			} else {
				_ = os.Setenv("GF_PLUGIN_BLOCK_LIST", tt.blockList)
			}

			result := IsPathBlocked(tt.path)

			if result != tt.shouldBlock {
				t.Errorf("Expected IsPathBlocked to return %v, but got %v", tt.shouldBlock, result)
			}
		})
	}
}

func TestIsPathBlocked_DisableSecurityBlocklistEnv(t *testing.T) {
	// Save and restore the original env var
	originalValue := os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST")
	defer func() {
		if originalValue == "" {
			_ = os.Unsetenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST")
		} else {
			_ = os.Setenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST", originalValue)
		}
	}()

	if !IsPathBlocked(defaultSecurityBlockList[0]) {
		t.Errorf("Expected IsPathBlocked to return true when GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST is unset, but got false")
	}

	_ = os.Setenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST", "true")

	// Should NOT be blocked when env var is set
	if IsPathBlocked(defaultSecurityBlockList[0]) {
		t.Errorf("Expected IsPathBlocked to return false when GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST is 'true', but got true")
	}
}

func TestIsPathBlocked_DisableGrafanaBlocklistEnv(t *testing.T) {
	// Save and restore the original env var
	originalValue := os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST")
	defer func() {
		if originalValue == "" {
			_ = os.Unsetenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST")
		} else {
			_ = os.Setenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST", originalValue)
		}
	}()

	if !IsPathBlocked(defaultGrafanaInternalBlockList[0]) {
		t.Errorf("Expected IsPathBlocked to return true when GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST is unset, but got false")
	}

	_ = os.Setenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST", "true")

	// Should NOT be blocked when env var is set
	if IsPathBlocked(defaultGrafanaInternalBlockList[0]) {
		t.Errorf("Expected IsPathBlocked to return false when GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST is 'true', but got true")
	}
}
