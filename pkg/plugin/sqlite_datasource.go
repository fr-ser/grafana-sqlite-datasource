package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	// register sqlite driver
	_ "modernc.org/sqlite"
)

// Make sure sqliteDatasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime.
var (
	_ backend.QueryDataHandler   = (*sqliteDatasource)(nil)
	_ backend.CheckHealthHandler = (*sqliteDatasource)(nil)
)

// sqliteDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type sqliteDatasource struct {
	pluginConfig pluginConfig
}

type pluginConfig struct {
	Path        string
	PathOptions string
	PathPrefix  string
}

// NewDataSource creates a new datasource instance.
func NewDataSource(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Info("Creating instance")
	var config pluginConfig

	err := json.Unmarshal(setting.JSONData, &config)
	if err != nil {
		log.DefaultLogger.Error("Could unmarshal settings of data source", "err", err)
		return &sqliteDatasource{}, fmt.Errorf("error while unmarshalling data source settings: %s", err)
	}

	securePathOptions, securePathOptionsExist := setting.DecryptedSecureJSONData["securePathOptions"]
	if securePathOptionsExist {
		if config.PathOptions == "" {
			config.PathOptions = securePathOptions
		} else {
			config.PathOptions = config.PathOptions + "&" + securePathOptions
		}
	}

	return &sqliteDatasource{pluginConfig: config}, nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (ds *sqliteDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (
	*backend.QueryDataResponse, error,
) {
	log.DefaultLogger.Debug("Received request for data")
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		response.Responses[q.RefID] = query(q, ds.pluginConfig)
		log.DefaultLogger.Debug("Finished query", "refID", q.RefID)
	}

	return response, nil
}
