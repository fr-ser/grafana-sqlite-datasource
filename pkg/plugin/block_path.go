package plugin

import (
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

func IsPathBlocked(path string) bool {
	// Normalize path to lowercase for case-insensitive matching
	lowerPath := strings.ToLower(path)

	if os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST") != "true" {
		for _, term := range defaultSecurityBlockList {
			if strings.Contains(lowerPath, strings.ToLower(term)) {
				return true
			}
		}
	}

	if os.Getenv("GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST") != "true" {
		for _, term := range defaultGrafanaInternalBlockList {
			if strings.Contains(lowerPath, strings.ToLower(term)) {
				return true
			}
		}
	}

	// Check against additional user-defined blocklist from environment variable
	blockList, exists := os.LookupEnv("GF_PLUGIN_BLOCK_LIST")
	if exists && blockList != "" {
		blockedTerms := strings.Split(blockList, ",")
		for _, term := range blockedTerms {
			term = strings.TrimSpace(term)
			if term != "" && strings.Contains(lowerPath, strings.ToLower(term)) {
				return true
			}
		}
	}

	return false
}
