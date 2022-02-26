package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func checkDB(pathPrefix string, path string, options string) error {
	if pathPrefix == "file:" || pathPrefix == "" {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no file exists at the file path")
		} else if err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite", pathPrefix+path+"?"+options)
	if err != nil {
		return err
	}
	defer db.Close()

	sth, err := db.Exec("pragma schema_version;")
	if err != nil {
		return err
	}
	fmt.Println(sth)

	return nil
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

	err = checkDB(config.PathPrefix, config.Path, config.PathOptions)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("error checking db: %s", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
