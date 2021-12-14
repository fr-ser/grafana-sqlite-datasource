# Grafana SQLite Datasource

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)
[![CI](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/ci_cd.yml/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/ci_cd.yml)

This is a Grafana backend plugin to allow using an SQLite database as a data source.
The SQLite database needs to be accessible to the filesystem of the device where Grafana itself
is running.

## Plugin Installation

The most up to date (but also most generic) information can always be found here:
[Grafana Website - Plugin Installation](https://grafana.com/docs/grafana/latest/plugins/installation/#install-grafana-plugins)

### Recommended: Installing the Official and Released Plugin on an Existing Grafana With the CLI

Grafana comes with a command line tool that can be used to install plugins.

1. Run this command: `grafana-cli plugins install frser-sqlite-datasource`
2. Restart the Grafana server.
3. Login in with a user that has admin rights. This is needed to create datasources.
4. To make sure the plugin was installed, check the list of installed datasources. Click the
   Plugins item in the main menu. Both core datasources and installed datasources will appear.

### Latest Version: Installing the newest Plugin Version on an Existing Grafana With the CLI

The grafana-cli can also install plugins from a non-standard URL. This way even plugin versions,
that are not (yet) released to the official Grafana repository can be installed.

1. Run this command:

   ```sh
   # replace the $VERSION part in the URL below with the desired version (e.g. 2.0.2)
   grafana-cli --pluginUrl https://github.com/fr-ser/grafana-sqlite-datasource/releases/download/v$VERSION/frser-sqlite-datasource-$VERSION.zip plugins install frser-sqlite-datasource
   ```

2. See the recommended installation above (from the restart step)

### Manual: Installing the Plugin Manually on an Existing Grafana

In case the grafana-cli does not work for whatever reason plugins can also be installed manually.

1. Get the zip file from [Latest release on Github](https://github.com/fr-ser/grafana-sqlite-datasource/releases/latest)
2. Extract the zip file into the data/plugins subdirectory for Grafana:
   `unzip <the_download_zip_file> -d <plugin_dir>/`

   Finding the plugin directory can sometimes be a challenge as this is platform and settings
   dependent. A common location for this on Linux devices is `/var/lib/grafana/plugins/`
3. See the recommended installation above (from the restart step)

## Configuring the Datasource in Grafana

The only required configuration is the path to the SQLite database (local path on the Grafana Server)

1. Add an SQLite datasource.
2. Set the path to the database (the grafana process needs to find the SQLite database under this path).
3. Save the datasource and use it.

## Support for Time Formatted Columns

SQLite has no native "time" format. It relies on strings and numbers for time and dates. Since
especially for time series Grafana expects an actual time type, however, the plugin provides a way
to infer a real timestamp. This can be set in the query editor by providing the name of the column,
which should be reformatted to a timestamp.

The plugin supports two different inputs that can be converted to a "time" depending on the type
of the value in the column, that should be formatted as "time":

1. **A number input**: It is assumed to be a unix timestamp / unix epoch. This represents time in
   the number of **seconds** (make sure your timestamp is not in milliseconds). More information is
   here: <https://en.wikipedia.org/wiki/Unix_time>

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

## Macros

This plugins supports macros inspired by the built-in Grafana datasources (e.g.
<https://grafana.com/docs/grafana/latest/datasources/postgres/#macros>).

However, as each macro needs to be re-implemented from scratch, only the following macros are
supported. Other macros (that you might expect from other SQL databases) do not work

### $__unixEpochGroupSeconds(unixEpochColumnName, intervalInSeconds)

Example: `$__unixEpochGroupSeconds("time", 10)`

Will be replaced by an expression usable in GROUP BY clause. For example:
`cast(("time" / 10) as int) * 10`

### $__unixEpochGroupSeconds(unixEpochColumnName, intervalInSeconds, NULL)

Example: `$__unixEpochGroupSeconds(timestamp, 10, NULL)`

Same as above but with a fill parameter so missing points in that series will be added for Grafana
and `NULL` will be used as value.

In case multiple time columns are provided the first one is chosen as the column to determine the
gap filling. "First" in this context means first in the SELECT statement. This column needs to have
no NULL values and must sorted in ascending order.

## Alerting

This plugins supports the Grafana alerting feature. Similar to the built in data sources alerting
does not support variables as they are normally replaced in the frontend, which is not involved
for the alerts. In order to allow time filtering this plugin supports the variables `$__from` and
`$__to`. For more information about those variables see here:
<https://grafana.com/docs/grafana/latest/variables/variable-types/global-variables/#__from-and-__to>.
Formatting of those variables (e.g. `${__from:date:iso}`) is not supported for alerts, however.

### ARMv7 / Raspberry Pi 2 Mod. B Support

Many programs for ARMv6 (the default 32bit ARM version of this plugin) can run on ARMv7 machines
(the version running on Raspberry Pi 2 Mod. B).
In case there are complications you can try out the version specifically built for ARMv7, though.

A plugin version specifically built for ARMv7 devices can be found on the Github release page (see
manual installation above).

## Development and Contributing

Any contribution is welcome. Some information regarding the local setup can be found in the
[DEVELOPMENT.md file](https://github.com/fr-ser/grafana-sqlite-datasource/blob/master/DEVELOPMENT.md).
