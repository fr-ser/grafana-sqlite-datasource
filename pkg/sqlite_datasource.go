package main

import (
	"context"
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &SQLiteDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

type pluginConfig struct {
	Path string
}

func getConfig(settings *backend.DataSourceInstanceSettings) (pluginConfig, error) {
	var config pluginConfig
	err := json.Unmarshal(settings.JSONData, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// SQLiteDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type SQLiteDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (td *SQLiteDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (
	*backend.QueryDataResponse, error,
) {
	log.DefaultLogger.Info("QueryData", "req", req)
	response := backend.NewQueryDataResponse()

	config, err := getConfig(req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		log.DefaultLogger.Error("Could not get config for plugin", "err", err)
		return response, err
	}

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		// the response can have an error field
		response.Responses[q.RefID] = query(q, config)
	}

	return response, nil
}

type instanceSettings struct{}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (
	instancemgmt.Instance, error,
) {
	log.DefaultLogger.Info("Creating instance")
	return &instanceSettings{}, nil
}

// Dispose is called before creating a new instance to allow plugin authors to cleanup
func (s *instanceSettings) Dispose() {
	log.DefaultLogger.Info("Disposing of instance")
}
