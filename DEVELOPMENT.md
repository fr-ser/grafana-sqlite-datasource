# Grafana SQLite Datasource

This is a Grafana backend plugin to allow using a SQLite database as a data source.

The plugin was built using the grafana plugin sdk and npx grafana toolkit. Information can be
found at:

- <https://grafana.com/tutorials/build-a-data-source-backend-plugin/>
- <https://github.com/grafana/grafana-plugin-sdk-go>
- <https://github.com/grafana/grafana/tree/main/packages/grafana-toolkit>

## Getting started

This project uses the `Makefile` as the main script tool for development tasks. Run `make help` to
get an overview of the available commands.

### Requirements

- nodejs
- go
- docker and docker compose
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

Now you can connect to the dockerized browser via a `VNC` client/viewer.

The VNC password is `secret`.

**Note - Macos M1/AMD support:** The selenium image does not support arm architectures yet.

In order to run the tests on an ARM architecture please change the docker image of selenium by using this environment variable:
`SELENIUM_IMAGE=seleniarm/standalone-chromium:112.0`.

You can find more information here: <https://github.com/SeleniumHQ/docker-selenium#experimental-mult-arch-aarch64armhfamd64-images>

#### VNC Viewer

On linux distributions you can use remmina as a VNC viewer.

On MacOs you can use the preinstalled "screen sharing" application as a VNC viewer.

## Release process

After step 3 Github Actions should take over and create a new release.
Steps 4 and 5 are for publishing the release to Grafana repository.

1. Make sure a section in the Changelog exists with `## [Unreleased]`
2. Push the changes and merge to the default branch
3. Get the md5 hash of the release from the Github Action or from the release page (text file)
4. Within the Grafana Cloud account a request for a plugin update can be started:
   <https://grafana.com/orgs/frser/plugins>
