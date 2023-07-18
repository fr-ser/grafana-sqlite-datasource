# Grafana SQLite Datasource

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)
[![branches](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/branches.yml/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/branches.yml)
[![tags](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/tags.yml/badge.svg)](https://github.com/fr-ser/grafana-sqlite-datasource/actions/workflows/tags.yml)

This is a Grafana backend plugin to allow using an SQLite database as a data source.
The SQLite database needs to be accessible to the filesystem of the device where Grafana itself
is running.

Table of contents:

- [Plugin Installation](#plugin-installation)
- [Support for Time Formatted Columns](#support-for-time-formatted-columns)
- [Macros](#macros)
- [Alerting](#alerting)
- [Common Problems - FAQ](#common-problems---faq)
- [Development and Contributing](#development-and-contributing)
- [Supporting the Project](#supporting-the-project)
- [Further Documentation and Links](#further-documentation-and-links)

## Plugin Installation

The recommended way to install the plugin for most users is to use the grafana CLI:

1. Run this command: `grafana-cli plugins install frser-sqlite-datasource`
2. Restart the Grafana server.
3. To make sure the plugin was installed, check the list of installed data sources. Click the
   Plugins item in the main menu. Both core data sources and installed data sources will appear.

For other installation options (e.g. to install versions not yet releases in the Grafana registry but in Github) see
[./docs/installation.md](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs/installation.md).

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

This plugins supports macros inspired by the built-in Grafana data sources (e.g.
<https://grafana.com/docs/grafana/latest/datasources/postgres/#macros>).

However, as each macro needs to be re-implemented from scratch, only the following macros are
supported. Other macros (that you might expect from other SQL databases) are not supported by the
plugin (yet).

### $__unixEpochGroupSeconds(unixEpochColumnName, intervalInSeconds)

Example: `$__unixEpochGroupSeconds("time", 10)`

Will be replaced by an expression usable in GROUP BY clause. For example:
`cast(("time" / 10) as int) * 10`

### $__unixEpochGroupSeconds(unixEpochColumnName, intervalInSeconds, NULL)

Example: `$__unixEpochGroupSeconds(timestamp, 10, NULL)`

This is the same as the above example but with a fill parameter so missing points in that series
will be added for Grafana and `NULL` will be used as value.

In case multiple time columns are provided the first one is chosen as the column to determine the
gap filling. "First" in this context means first in the SELECT statement. This column needs to have
no NULL values and must be sorted in ascending order.

## Alerting

The plugins supports the Grafana alerting feature. Similar to the built in data sources alerting
does not support variables as they are normally replaced in the frontend, which is not involved
for the alerts. In order to allow time filtering this plugin supports the variables `$__from` and
`$__to`. For more information about those variables see here:
<https://grafana.com/docs/grafana/latest/variables/variable-types/global-variables/#__from-and-__to>.
Formatting of those variables (e.g. `${__from:date:iso}`) is not supported for alerts, however.

## Common Problems - FAQ

This is a list of common questions or problems. For the answers and more details see
[./docs/faq.md](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs/faq.md).

- [I have a "file not found" error for my database](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs/faq.md#i-have-a-file-not-found-error-for-my-database)
- [I have a "permission denied" error for my database](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs/faq.md#i-have-a-permission-denied-error-for-my-database)
- ...

## Query examples

Some examples to help getting started with SQL and SQLite can be found in
[./docs/examples.md](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs/examples.md).

These examples include things like:

- filtering by the time specified in the Grafana UI
- creating a time series mindful of gaps
- converting dates to timestamps
- ...

## Development and Contributing

Any contribution is welcome. Some information regarding the local setup can be found in the
[DEVELOPMENT.md file](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/DEVELOPMENT.md).

## Supporting the Project

This project was developed for free as an open source project. And it will stay that way.

If you like using this plugin, however, and would like to support the development go check out
the [Github sponsorship page](https://github.com/sponsors/fr-ser). This allows sponsoring the
project with monthly or one-time contributions.

## Further Documentation and Links

- A changelog of the plugin can be found in the [CHANGELOG.md](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/CHANGELOG.md).
- More documentation about the plugin can be found under [the docs section in Github](https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/docs).
- The plugin in the Grafana registry can be found [here](https://grafana.com/grafana/plugins/frser-sqlite-datasource/).
- Questions or bugs about the plugin can be found and reported [in Github](https://github.com/fr-ser/grafana-sqlite-datasource/issues?q=) or in the [Grafana community](https://community.grafana.com/search?q=sqlite%20order%3Alatest).
