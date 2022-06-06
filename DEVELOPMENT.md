# Grafana SQLite Datasource

This is a Grafana backend plugin to allow using a SQLite database as a data source.

The plugin was built using the grafana plugin sdk and npx grafana toolkit. Information can be
found at:

- <https://grafana.com/tutorials/build-a-data-source-backend-plugin/>
- <https://github.com/grafana/grafana-plugin-sdk-go>
- <https://github.com/grafana/grafana/tree/master/packages/grafana-toolkit>

## Getting started

This project uses the `Makefile` as the main script tool for development tasks. Run `make help` to
get an overview of the available commands.

### Requirements

- nodejs
- yarn
- go
- docker and docker-compose
- make

### (First Time) Installation

```sh
# installing packages
make install
# optional: using git hooks
git config core.hooksPath githooks
```

### Start up Grafana

```sh
make build # this build the frontend and backend
make start # credentials admin / admin123
```

## Testing

```sh
make test
```

### Quick e2e tests with Selenium

First start the docker environment with `make selenium-test`. This will also run the tests.
Regardless of the tests passing the environment will stay up and running.

Now you can connect to the dockerized browser via a `VNC` client/viewer (like remmina)

## Release process

After step 3 Github Actions should take over and create a new release.
Steps 4 and 5 are for publishing the release to Grafana repository.

1. Update the Changelog
2. Tag the commit with a Semver tag, e.g. v2.2.3-rc.1
3. Push the changes including the tag
4. Get the md5 hash of the release from the Github Action or from the release page (text file)
5. Within the Grafana Cloud account a request for a plugin update can be started:
   <https://grafana.com/orgs/frser/plugins>
