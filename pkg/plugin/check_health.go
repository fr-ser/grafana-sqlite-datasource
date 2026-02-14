package plugin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Default blocklist of sensitive paths that should never be accessible
var defaultBlockList = []string{
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

	// Grafana internal
	"grafana.db",

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

func IsPathBlocked(path string) bool {
	// Normalize path to lowercase for case-insensitive matching
	lowerPath := strings.ToLower(path)

	// Check against default blocklist
	for _, term := range defaultBlockList {
		if strings.Contains(lowerPath, strings.ToLower(term)) {
			return true
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

func checkDB(pathPrefix string, path string, options string) error {
	if IsPathBlocked(path) {
		return fmt.Errorf("path contains blocked term from GF_PLUGIN_BLOCK_LIST")
	}

	if pathPrefix == "file:" || pathPrefix == "" {
		fileInfo, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no file exists at the file path")
		} else if err != nil {
			return fmt.Errorf("error checking path: %v", err)
		} else if fileInfo.IsDir() {
			return fmt.Errorf("the provided path is a directory instead of a file")
		}
	}

	db, err := sql.Open("sqlite", pathPrefix+path+"?"+options)
	if err != nil {
		return fmt.Errorf("error opening %s%s: %v", pathPrefix, path, err)
	}

	_, err = db.Exec("pragma schema_version;")
	if err != nil {
		return fmt.Errorf("error checking for valid SQLite file: %v", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing database file: %v", err)
	}

	return nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *sqliteDatasource) CheckHealth(ctx context.Context, _ *backend.CheckHealthRequest) (
	*backend.CheckHealthResult, error,
) {
	err := checkDB(ds.pluginConfig.PathPrefix, ds.pluginConfig.Path, ds.pluginConfig.PathOptions)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("error checking db: %s", err),
		}, nil
	}

	if ds.pluginConfig.AttachLimit != nil && *ds.pluginConfig.AttachLimit > 0 {
		if os.Getenv("GF_PLUGIN_UNSAFE_ALLOW_ATTACH_LIMIT_ABOVE_ZERO") != "true" {
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "An 'attach limit' above 0 is not allowed without setting 'unsafe_allow_attach_limit_above_zero' in the plugin configuration",
			}, nil
		}
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
