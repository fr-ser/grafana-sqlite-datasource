package plugin

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var ctx = context.Background()

func TestCheckHealthShouldPassForADB(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "my.db")

	db, _ := sql.Open("sqlite", dbPath)
	_, _ = db.Exec("CREATE TABLE test(id int);")
	db.Close()

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

func TestCheckHealthShouldFailIfNoFileExists(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	defer os.RemoveAll(dir)
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

func TestCheckHealthShouldFailOnTextFile(t *testing.T) {
	f, _ := ioutil.TempFile("", "test-check-db")
	defer func() { _ = syscall.Unlink(f.Name()) }()
	_, _ = f.WriteString("not a sqlite db")
	f.Close()

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
	f, _ := ioutil.TempFile("", "test-check-db")
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
