package plugin

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var ctx = context.Background()

func TestCheckHealthShouldPassForADB(t *testing.T) {
	dir, _ := os.MkdirTemp("", "test-check-db")
	defer func() { _ = os.RemoveAll(dir) }()
	dbPath := filepath.Join(dir, "my.db")

	db, _ := sql.Open("sqlite", dbPath)
	_, _ = db.Exec("CREATE TABLE test(id int);")
	_ = db.Close()

	ds := sqliteDatasource{pluginConfig{Path: dbPath, PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusOk {
		t.Errorf("Expected HealthStatusOk, but got - %s", result.Status)
	}
	if result.Message != "Data source is working" {
		t.Errorf("Unexpected message: %s", result.Message)
	}
}

func TestCheckHealthShouldFailForANotExistingFile(t *testing.T) {
	dir, _ := os.MkdirTemp("", "test-check-db")
	defer func() { _ = os.RemoveAll(dir) }()
	notExistingDbPath := filepath.Join(dir, "my.db")

	ds := sqliteDatasource{pluginConfig{Path: notExistingDbPath, PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}
	if !strings.Contains(result.Message, "no file exists at the file path") {
		t.Errorf("Unexpected error message: %s", result.Message)
	}

	_, err = os.Stat(notExistingDbPath)
	if !os.IsNotExist(err) {
		t.Errorf("File was created during check")
	}
}

func TestCheckHealthShouldFailForAFolder(t *testing.T) {
	dir, _ := os.MkdirTemp("", "test-check-db")
	defer func() { _ = os.RemoveAll(dir) }()

	ds := sqliteDatasource{pluginConfig{Path: dir, PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}

	if !strings.Contains(result.Message, "the provided path is a directory instead of a file") {
		t.Errorf("Unexpected error message: %s", result.Message)
	}
}

func TestCheckHealthShouldFailOnTextFile(t *testing.T) {
	f, _ := os.CreateTemp("", "test-check-db")
	defer func() { _ = syscall.Unlink(f.Name()) }()
	_, _ = f.WriteString("not a sqlite db")
	_ = f.Close()

	ds := sqliteDatasource{pluginConfig{Path: f.Name(), PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}
	if !strings.Contains(result.Message, "file is not a database") {
		t.Errorf("Unexpected error message: %s", result.Message)
	}
}

func TestCheckHealthShouldPassForAnEmptyFile(t *testing.T) {
	f, _ := os.CreateTemp("", "test-check-db")
	defer func() { _ = syscall.Unlink(f.Name()) }()

	ds := sqliteDatasource{pluginConfig{Path: f.Name(), PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusOk {
		t.Errorf("Expected HealthStatusOk, but got - %s", result.Status)
	}
}

func TestCheckHealthShouldFailWhenPathIsBlocked(t *testing.T) {
	dbPath := "/some/root/path/secret-database.db"

	originalValue := os.Getenv("GF_PLUGIN_BLOCK_LIST")
	defer func() {
		if originalValue == "" {
			_ = os.Unsetenv("GF_PLUGIN_BLOCK_LIST")
		} else {
			_ = os.Setenv("GF_PLUGIN_BLOCK_LIST", originalValue)
		}
	}()
	_ = os.Setenv("GF_PLUGIN_BLOCK_LIST", "secret")

	ds := sqliteDatasource{pluginConfig{Path: dbPath, PathPrefix: "file:"}}
	result, err := ds.CheckHealth(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}
	if !strings.Contains(result.Message, "path contains blocked term from GF_PLUGIN_BLOCK_LIST") {
		t.Errorf("Unexpected error message: %s", result.Message)
	}
}

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
			name:        "case sensitive matching",
			blockList:   "Secret",
			path:        "/some/path/secret-database.db",
			shouldBlock: false,
		},
		{
			name:        "partial path matching",
			blockList:   "tmp",
			path:        "/tmp/database.db",
			shouldBlock: true,
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
