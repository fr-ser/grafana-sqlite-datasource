package plugin

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

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
