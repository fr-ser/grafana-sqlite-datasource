package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mattn/go-sqlite3"
)

func checkDB(pathPrefix string, path string, options string) (bool, error) {
	// to avoid creating a file during this check if it did not exist we try to append the
	// read only mode. We do not overwrite an existing mode setting, however.
	// if the pathPrefix is not "file:", this readonly mode setting has no effect
	var finalOptions string
	if options == "" {
		finalOptions = "mode=ro"
	} else if !strings.Contains(options, "mode") {
		finalOptions = finalOptions + "&mode=ro"
	} else {
		finalOptions = options
	}

	db, err := sql.Open("sqlite3", pathPrefix+path+"?"+finalOptions)
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

	isDB, err := checkDB(config.PathPrefix, config.Path, config.PathOptions)
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
