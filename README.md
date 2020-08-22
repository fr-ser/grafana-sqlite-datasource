# Grafana SQLite Datasource

This is a Grafana backend plugin to allow using a SQLite database as a data source.

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
