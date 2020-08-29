# Grafana SQLite Datasource

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![stability-wip](https://img.shields.io/badge/stability-work_in_progress-lightgrey.svg)

[![CI Tests](https://github.com/fr-ser/grafana-sqlite-datasource/workflows/Test%20%26%20Build/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions)

This is a Grafana backend plugin to allow using a SQLite database as a data source.

The plugin was built using the grafana plugin sdk and npx grafana toolkit. Information can be
found at:

- https://grafana.com/tutorials/build-a-data-source-backend-plugin/
- https://github.com/grafana/grafana-plugin-sdk-go
- https://github.com/grafana/grafana/tree/master/packages/grafana-toolkit

## Getting started

### Requirements

- yarn
- go
- docker-compose

### (First Time) Installation

```BASH
make install
```

### Start up Grafana

```BASH
make build # this build the frontend and backend
mage bootstrap # credentials admin / admin123
```

## Testing

```BASH
make test
```

## TODO: Cross compilation

Resources

- https://www.arp242.net/static-go.html
- https://dh1tw.de/2019/12/cross-compiling-golang-cgo-projects/
- https://github.com/grafana/google-sheets-datasource/issues/104
- https://github.com/grafana/grafana-plugin-sdk-go/issues/188
