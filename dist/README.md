# Grafana SQLite Datasource

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![stability-wip](https://img.shields.io/badge/stability-work_in_progress-lightgrey.svg)

[![CI Tests](https://github.com/fr-ser/grafana-sqlite-datasource/workflows/Test%20%26%20Build/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions)

This is a Grafana backend plugin to allow using a SQLite database as a data source.

## Development and Contributing

Any contribution is welcome. Some information regarding the local setup can be found in the
[DEVELOPMENT.md file](https://github.com/fr-ser/grafana-sqlite-datasource/blob/master/DEVELOPMENT.md).

## Plugin installation

The most up to date (but also most generic) information can always be found here:
[Grafana Website - Plugin Installation](https://grafana.com/docs/grafana/latest/plugins/installation/#install-grafana-plugins)

### Installing the Plugin on an Existing Grafana with the CLI

Grafana comes with a command line tool that can be used to install plugins.

1. Run this command: `grafana-cli plugins install frser-sqlite-datasource`
2. Restart the Grafana server.
3. Login in with a user that has admin rights. This is needed to create datasources.
4. To make sure the plugin was installed, check the list of installed datasources. Click the Plugins item in the main menu. Both core datasources and installed datasources will appear.

### Installing the Plugin Manually on an Existing Grafana

If the server where Grafana is installed has no access to the Grafana.com server, then the plugin can be downloaded and manually copied to the server.

2. Get the zip file from https://github.com/fr-ser/grafana-sqlite-datasource/archive/v0.1.2.zip
3. Extract the zip file into the data/plugins subdirectory for Grafana:
   `unzip grafana-sqlite-datasource-0.1.2.zip -d YOUR_PLUGIN_DIR/grafana-sqlite-datasource`
4. Restart the Grafana server
5. To make sure the plugin was installed, check the list of installed datasources. Click the Plugins item in the main menu. Both core datasources and installed datasources will appear.

## Configuring the datasource in Grafana

The only required configuration is the path to the SQLite database (local path on the Grafana Server)

1. Add an SQLite datasource.
2. Set the path to the database (the grafana process needs to find the SQLite database under this path).
3. Save the datasource and use it.