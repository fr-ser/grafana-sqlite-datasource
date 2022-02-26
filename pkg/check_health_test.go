package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var ctx = context.Background()

func getReqWithPath(path string) *backend.CheckHealthRequest {
	jsonConfig := fmt.Sprintf(`{"pathPrefix":"file:", "path": "%s"}`, path)

	return &backend.CheckHealthRequest{
		PluginContext: backend.PluginContext{
			DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
				JSONData: []byte(jsonConfig),
			},
		},
	}
}

func TestCheckHealthShouldPassForADB(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "my.db")

	db, _ := sql.Open("sqlite", dbPath)
	db.Exec("CREATE TABLE test(id int);")
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

func TestCheckHealthShouldFailIfNoFileExists(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	defer os.RemoveAll(dir)
	notExistingDbPath := filepath.Join(dir, "my.db")

	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath(notExistingDbPath))
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
	defer syscall.Unlink(f.Name())
	f.WriteString("not a sqlite db")
	f.Close()

	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath(f.Name()))
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
	defer syscall.Unlink(f.Name())

	ds := SQLiteDatasource{}
	result, err := ds.CheckHealth(ctx, getReqWithPath(f.Name()))
	if err != nil {
		t.Errorf("Unexpected error - %s", err)
	}

	if result.Status != backend.HealthStatusOk {
		t.Errorf("Expected HealthStatusOk, but got - %s", result.Status)
	}
}
