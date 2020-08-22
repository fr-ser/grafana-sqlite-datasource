package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var ctx = context.Background()

func getReqWithPath(path string) *backend.CheckHealthRequest {
	jsonConfig := fmt.Sprintf(`{"path": "%s"}`, path)

	return &backend.CheckHealthRequest{
		PluginContext: backend.PluginContext{
			DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
				JSONData: []byte(jsonConfig),
			},
		},
	}
}

func TestCheckAllGood(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "my.db")

	db, _ := sql.Open("sqlite3", dbPath)
	db.Exec("SELECT 1")
	db.Close()

	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath(dbPath))
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

func TestCheckFileExistence(t *testing.T) {
	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath("hello"))
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}
	if result.Message != "No file found at: 'hello'" {
		t.Errorf("Unexpected error message: %s", result.Message)
	}
}

func TestCheckIsDatabase(t *testing.T) {
	f, _ := ioutil.TempFile("", "test-check-db")
	defer syscall.Unlink(f.Name())
	f.WriteString("not a sqlite db")

	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath(f.Name()))
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("Expected HealthStatusError, but got - %s", result.Status)
	}
	if result.Message != "The file at the provided location, was not a valid SQLite database" {
		t.Errorf("Unexpected error message: %s", result.Message)
	}
}
