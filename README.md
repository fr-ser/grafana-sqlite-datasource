# Grafana SQLite Datasource

![stability-wip](https://img.shields.io/badge/stability-work_in_progress-lightgrey.svg)

![CI Tests](https://github.com/fr-ser/grafana-sqlite-datasource/workflows/main/badge.svg)

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
- mage: Makefile with go - comes with the plugin toolkit 🤷
- docker-compose

### (First Time) Installation

```BASH
mage -v install
```

### Start up Grafana

```BASH
mage -v # this build the frontend and backend
mage bootstrap # credentials admin / admin123
```

## Testing

Currently there are only backend (go) tests. Run via:

```BASH
mage -v test
```
