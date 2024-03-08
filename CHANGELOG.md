# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com) and follows semantic Versioning.

## [3.5.0] - 2024-03-08

### Added

- Added the ability to use "__text" and "__value" columns for query variables

## [3.4.0] - 2023-12-26

Thank you to the collaborators on this release:

- [@eeddaann](https://github.com/eeddaann)

### Added

- added support for Grafana annotations

## [3.3.3] - 2023-07-21

### Changed

- updated frontend packages with the grafana/create-plugin tool

## [3.3.2] - 2023-07-19

### Fixed

- queries without any content (i.e. "-- not a query") do not cause the plugin to hang

### Changed

- updated the underlying sqlite library

## [3.3.1] - 2023-07-01

### Fixed

- Closing database connections after every query (bug from v3.3.0)

### Changed

- Functional regression: a query without any content (i.e. "-- not a query") causes the plugin to hang

## [3.3.0] - 2023-06-11

### Changed

- updated sqlite package version to allow setting limits
- add option to set `SQLITE_LIMIT_ATTACHED` (current default 10).
  In a future version the default will be reduced to 0.

## [3.2.0] - 2023-05-25

### Changed

- updated go version and go grafana plugin sdk
- use grafana plugin sdk mage file to build the plugin
- update grafana plugin npm packages

## [3.1.4] - 2023-05-17

Nothing user facing was changed.
The only difference is a change in meta data to accommodate new Grafana plugin checks.

## [3.1.3] - 2023-05-06

### Changed

- updated go SQLite package to v1.22.1 to fix timezone issues in SQLite

## [3.1.2] - 2022-09-11

### Changed

- updated go SQLite package

## [3.1.1] - 2022-09-11

### Changed

- removed the "name" metadata from data frames

## [3.1.0] - 2022-07-09

### Changed

- Using a different component as the SQL editor
- Updated Readme and added more docs (e.g. example queries)
- improved error message for database checks on folder paths

I would like to thank the sponsors of this project

- one-time: [@nikolaialex](https://github.com/nikolaialex)

## [3.0.1] - 2022-03-15

This release should show no feature changes.
Some underlying packages have been updated, which should improve stability and security but not have any other noticeable impact.

### Changed

- Fixed some typos in the readme
- Removed the upper Grafana version constraint for the plugin
- Update grafana plugin sdk for the backend
- Updated grafana buildkit and test tool versions for the frontend

## [3.0.0] - 2022-02-26

This release moved to a new underlying SQLite library: <https://pkg.go.dev/modernc.org/sqlite>.
This should have no big changes to regular queries but can have effects on more subtle configurations (e.g. path options).
Fore more information see the `Changed` section below.

This library has no dependency on CGO, which allows much easier cross-compilation for other systems. This way the
plugin has a much simpler build process now and also supports more platforms (see information below under `Added`)

### Added

- Added new platforms to the release: Darwin (MacOS) for ARM (Apple Silicon) and FreeBSD for AMD64

### Changed

- Changed the underlying SQLite library to: <https://pkg.go.dev/modernc.org/sqlite>.
  While the general SQLite features and especially queries should remain unchanged by this, path options need to be checked for compatibility with the new library now.
  Please refer to the link above for more information on the options.

### Removed

- No separate release is created anymore for ARM v7. ARM v6 should suffice for all devices. If that
  is not the case for your device, please open a new issue.

## [2.2.1] - 2021-12-14

### Added

- An additional option `securePathOptions` has been added in case the user wants to protect some options (typically credentials).
  For examples see here: <https://github.com/mattn/go-sqlite3#connection-string>

## [2.2.0] - 2021-11-16

The plugin now supports adding a Path Prefix and Options to the SQLite connection string.

### Added

- Ability to provide a prefix and options in the connection string. For examples see here:
  <https://github.com/mattn/go-sqlite3#connection-string>

### Changed

- slightly changed the Plugin Health Check (when adding the data source) to provide better error messages.
- conversion errors during a query are now logged at DEBUG level to avoid too large log files.

## [2.1.1] - 2021-10-24

This release adds support for sub second precision for unix time.

### Added

- When using numeric values for a timestamp in SQLite (unix timestamp) the plugin now supports precision below the second (at nanosecond precision)

## [2.1.0] - 2021-08-08

This release adds the JSON extension to the compiled SQLite code.

### Added

- JSON extension for SQLite

## [2.0.2] - 2021-08-07

This release adds testing against Grafana v8.1.0 and fixes an issue with query variables.

### Fixed

- Query variables can now also be used in Grafana v8.X.X

## [2.0.1] - 2021-07-27

This release fixes some long standing issues that prevented the right use of the alerting feature of the plugin even though it was enabled already.

### Fixed

- Using the `$__from` and `$__to` variables for alerting
- Fixing a caching bug for the query (for alerting)

## [2.0.0] - 2021-05-07

All current Raspberry Pi Zero and 1 models have an ARMv6 architecture.
All other models (2 Mod. B v1.2, 3 and 4) have an 64Bit ARMv8 architecture.
As only the Raspberry Pi 2 Mod. B has an ARMv7 architecture this is not used as default anymore.
The Raspberry Pi 2 Mod. B will require a manual installation and all others will be handled via the Grafana CLI.

### Changed

- Using ARMv6 instead of ARMv7 as 32Bit ARM default

## [1.2.1] - 2021-04-05

### Added

- More debug level logging from the plugin

### Fixed

- The type inference of columns in the backend is now ignoring the letter casing

## [1.2.0] - 2021-03-31

### Added

- The response of the plugin includes the final query as metadata and can be checked in the inspector now
- Macro `unixEpochGroupSeconds`:
  - replace time columns with an expression to group by
  - Allow filling up missing values with `NULL`

### Fixed

- return additional time formatted column for time-series formats as normal values (previously they were skipped)

## [1.1.0] - 2021-03-27

### Added

- Experimental support for MacOS (no static linking)

## [1.0.3] - 2021-03-26

### Fixed

- Showing better error messages for certain fail conditions of the plugin health check (e.g. permission error)

## [1.0.2] - 2021-03-23

### Fixed

- Fixed bug preventing using query variables when SQLite is the default datasource (<= Grafana 7.4)

## [1.0.1] - 2021-03-21

### Added

- Enabled the `alerting` feature for the plugin (no code change)

## [1.0.0] - 2021-03-13

No breaking change was introduced but due to code stability the first major version is released.

### Fixed

- variables like `$__interval` and `$__interval_ms` are supported now

## [0.2.7] - 2021-01-07

### Changed

- Changing plugin name to SQLite
- added category to plugin.json for better grouping on the Grafana homepage
- updated Readme after first official release of plugin on Grafana homepage

## [0.2.6] - 2021-01-03

### Added

- Documentation about the time conversion in the README and in the query editor help text.

## [0.2.5] - 2020-12-19

### Fixed

- Correct handling of "NUMERIC" columns with mixed data (e.g. float and integer)

## [0.2.4] - 2020-12-13

### Added

- Added option to explicitly convert backend data frame to time series:

  - This requires the initial data frame to be in a [Long Format](https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#long-format)
  - The resulting time series consists of one data frame per metric

## [0.2.3] - 2020-11-28

### Changed

- Releasing arm6 (RaspberryPi Zero) as separate distribution (Github only)

### Fixed

- Renamed the arm7 executable to arm (newer Raspberry Models should run fine now)

## [0.2.2] - 2020-12-24

### Changed

- Different content of zip file published with Github release according to new Grafana v7.3 standards

## [0.2.1] - 2020-11-22

### Added

- Query variables are now supported

## [0.2.0] - 2020-11-20

### Added

- The plugin is now signed

### Changed

- For Signing grafana-toolkit 7.3.3 was necessary. The grafana version to test against was bumped to version 7.3.3

## [0.1.3] - 2020-11-02

### Fixed

- Fixed: Handling of NULL values in queries is now correct

## [0.1.2] - 2020-10-28

First "working" version

### Fixed

- Fixed: Plugin files in the zip file are now executable
