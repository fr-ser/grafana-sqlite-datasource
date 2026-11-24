package plugin

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Default security related blocklist of sensitive paths that should never be accessible
var defaultSecurityBlockList = []string{
	// when updating this least remember to also update the readme/documentation.

	// Cloud provider credentials
	".aws",
	".config/gcloud",
	".azure",
	".kube/config",
	".docker/config",

	// SSH and crypto keys
	".ssh",
	".gnupg",
	".pki",

	// System sensitive files
	"/etc/shadow",
	"/etc/passwd",
	"/etc/gshadow",
	"/proc/",
	"/sys/",

	// Common secret file patterns
	".env",
	"credentials",
	".git/config",
	".netrc",
	".npmrc",
	".pypirc",

	// Private keys
	"id_rsa",
	"id_dsa",
	"id_ecdsa",
	"id_ed25519",
}

// Default grafana internal blocklist of sensitive paths that should never be accessible
var defaultGrafanaInternalBlockList = []string{
	// when updating this least remember to also update the readme/documentation.

	"grafana.db",
}

func IsPathBlocked(path string) (bool, error) {
	// Also check the percent-decoded path to catch URI-encoded bypasses like
	// "grafana%2Edb" which SQLite resolves to "grafana.db".
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return true, fmt.Errorf("invalid URL encoding in database path: %w", err)
	}

	lowerPath := strings.ToLower(path)
	lowerDecoded := strings.ToLower(decoded)

	contains := func(term string) bool {
		lowerTerm := strings.ToLower(term)
		return strings.Contains(lowerPath, lowerTerm) || strings.Contains(lowerDecoded, lowerTerm)
	}

	if os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST") != "true" {
		for _, term := range defaultSecurityBlockList {
			if contains(term) {
				return true, nil
			}
		}
	}

	if os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST") != "true" {
		for _, term := range defaultGrafanaInternalBlockList {
			if contains(term) {
				return true, nil
			}
		}
	}

	// Check against additional user-defined blocklist from environment variable
	blockList, exists := os.LookupEnv("GF_PLUGIN_BLOCK_LIST")
	if exists && blockList != "" {
		blockedTerms := strings.Split(blockList, ",")
		for _, term := range blockedTerms {
			term = strings.TrimSpace(term)
			if term != "" && contains(term) {
				return true, nil
			}
		}
	}

	return false, nil
}
