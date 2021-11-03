package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mattn/go-sqlite3"
)

func checkDbExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func checkIsDB(path string) (bool, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return false, err
	}
	defer db.Close()

	_, err = db.Exec("pragma schema_version;")
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrNotADB {
			return false, nil
		}
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (td *SQLiteDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (
	*backend.CheckHealthResult, error,
) {
	config, err := getConfig(req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("error getting config: %s", err),
		}, nil
	}

	dbExists, err := checkDbExists(config.Path)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("error checking existence: %s", err),
		}, nil
	} else if !dbExists {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("No file found at: '%s'", config.Path),
		}, nil
	}

	isDB, err := checkIsDB(config.Path)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("error checking db: %s", err),
		}, nil
	} else if !isDB {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "The file at the provided location is not a valid SQLite database",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
