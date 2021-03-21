# Grafana SQLite Datasource

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)
[![Test Branch](https://github.com/fr-ser/grafana-sqlite-datasource/workflows/Test%20Branch/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions)
[![Test Release](https://github.com/fr-ser/grafana-sqlite-datasource/workflows/Test%20Release/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions)

This is a Grafana backend plugin to allow using a SQLite database as a data source.

## Development and Contributing

Any contribution is welcome. Some information regarding the local setup can be found in the
[DEVELOPMENT.md file](https://github.com/fr-ser/grafana-sqlite-datasource/blob/master/DEVELOPMENT.md).

## Plugin Installation

The most up to date (but also most generic) information can always be found here:
[Grafana Website - Plugin Installation](https://grafana.com/docs/grafana/latest/plugins/installation/#install-grafana-plugins)

### Installing the Plugin on an Existing Grafana With the CLI

Grafana comes with a command line tool that can be used to install plugins.

1. Run this command: `grafana-cli plugins install frser-sqlite-datasource`
2. Restart the Grafana server.
3. Login in with a user that has admin rights. This is needed to create datasources.
4. To make sure the plugin was installed, check the list of installed datasources. Click the
   Plugins item in the main menu. Both core datasources and installed datasources will appear.

### Installing the Plugin Manually on an Existing Grafana (Most up to Date)

If you need a version that is not released (yet) on the Grafana homepage or if the server where
Grafana is installed has no access to the Grafana.com server, then the plugin can be downloaded
and manually copied to the server.

1. Get the zip file from [Latest release on Github](https://github.com/fr-ser/grafana-sqlite-datasource/releases/latest)
2. Extract the zip file into the data/plugins subdirectory for Grafana:
   `unzip <the_download_zip_file> -d <plugin_dir>/`

   Finding the plugin directory can sometimes be a challenge as this is platform and settings
   dependent. A common location for this on Linux devices is `/var/lib/grafana/plugins/`
3. Restart the Grafana server
4. To make sure the plugin was installed, check the list of installed datasources. Click the
   Plugins item in the main menu. Both core datasources and installed datasources will appear.

### ARM6 / RaspberryPi Zero W Support

This plugins supports ARM6 (the version running on RaspberryPi Zero W). There is a problem, though,
with Grafana supporting ARM7 (newer Raspberries) and ARM6 at the same time. Grafana chooses
the correct plugin by file name. But both ARM6 and ARM7 are named
`<plugin_dir>/frser-sqlite-datasource/gpx_sqlite-datasource_linux_arm`.

Currently the ARM7 build is named like this by default, which is why the "official" plugin
distribution does not support ARM6 devices.

A plugin version specifically built for ARM6 devices can be found on the Github release page (see
manual installation above).

## Configuring the Datasource in Grafana

The only required configuration is the path to the SQLite database (local path on the Grafana Server)

1. Add an SQLite datasource.
2. Set the path to the database (the grafana process needs to find the SQLite database under this path).
3. Save the datasource and use it.

## Support for Time Formatted Columns

SQLite has no native "time" format. It actually relies on strings and numbers. Since especially
for time series Grafana expects an actual time type, however, the plugin provides a way to infer
a real timestamp. This can be set in the query editor by providing the name of the column, which
should be reformatted to a timestamp.

The plugin supports two different inputs that can be converted to a "time" depending on the type
of the value in the column, that should be formatted as "time":

1. **A number input**: It is assumed to be a unix timestamp / unix epoch and will be converted to
   an integer before converting it to a timestamp.

2. **A string input**: The value is expected to be formatted in accordance with **RFC3339**,
   e.g. `"2006-01-02T15:04:05Z07:00"`. Edge cases might occur and the parsing library used is the
   source of truth here: <https://golang.org/pkg/time/#pkg-constants>.

Timestamps stored as unix epoch should work out of the box, but the string formatting might require
adjusting your current format. The below example shows how to convert a "date" column to a parsable
timestamp:

```SQL
WITH converted AS (
   -- a row looks like this (value, date): 1.45, '2020-12-12'
   SELECT value,  date || 'T00:00:00Z' AS datetime FROM raw_table
)
SELECT datetime, value FROM converted ORDER BY datetime ASC
```
