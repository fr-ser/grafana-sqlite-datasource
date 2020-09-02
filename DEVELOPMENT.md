# Grafana SQLite Datasource

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
make bootstrap # credentials admin / admin123
```

## Testing

```BASH
make test
```
