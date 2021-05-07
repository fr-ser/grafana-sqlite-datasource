# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [2.0.0-rc.1]

All current Raspberry PI Zero and 1 models have an ARMv6 architecture.
All other models (2 Mod. B v1.2, 3 and 4) have an 64Bit ARMv8 architecture.
As only the Raspberry Pi 2 Mod. B has an ARMv7 architecture this is not used as default anymore.
The Raspberry Pi 2 Mod. B will require a manual installation and all others will be handled
via the Grafana CLI.

### Changed

- Using ARMv6 instead of ARMv7 as 32Bit ARM default

## [1.2.1]

### Added

- More debug level logging from the plugin

### Fixed

- The type inference of columns in the backend is now ignoring the letter casing

## [1.2.0]

### Added

- The response of the plugin includes the final query as metadata and can be checked in the
  inspector now
- Macro `unixEpochGroupSeconds`:
  - replace time columns with an expression to group by
  - Allow filling up missing values with `NULL`

### Fixed

- return additional time formatted column for time-series formats as normal values (previously
  they were skipped)

## [1.1.0]

### Added

- Experimental support for MacOS (no static linking)

## [1.0.3]

### Fixed

- Showing better error messages for certain fail conditions of the plugin health check (e.g.
  permission error)

## [1.0.2]

### Fixed

- Fixed bug preventing using query variables when SQLite is the default datasource (<= Grafana 7.4)

## [1.0.1]

### Added

- Enabled the `alerting` feature for the plugin (no code change)

## [1.0.0]

No breaking change was introduced but due to code stability the first major version is released.

### Fixed

- variables like `$__interval` and `$__interval_ms` are supported now

## [0.2.7]

### Changed

- Changing plugin name to SQLite

- added category to plugin.json for better grouping on the Grafana homepage

- updated Readme after first official release of plugin on Grafana homepage

## [0.2.6]

### Added

- Documentation about the time conversion in the README and in the query editor help text.

## [0.2.5]

### Fixed

- Correct handling of "NUMERIC" columns with mixed data (e.g. float and integer)

## [0.2.4]

### Added

- Added option to explicitly convert backend data frame to time series:

  - This requires the initial data frame to be in a [Long Format](https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#long-format)

  - The resulting time series consists of one data frame per metric

## [0.2.3]

### Changed

- Releasing arm6 (RaspberryPi Zero) as separate distribution (Github only)

### Fixed

- Renamed the arm7 executable to arm (newer Raspberry Models should run fine now)

## [0.2.2]

### Changed

- Different content of zip file published with Github release according to new Grafana v7.3
  standards

## [0.2.1]

### Added

- Query variables are now supported

## [0.2.0]

### Added

- The plugin is now signed

### Changed

- For Signing grafana-toolkit 7.3.3 was necessary. The grafana version to test against was
  bumped to version 7.3.3

## [0.1.3]

### Fixed

- Fixed: Handling of NULL values in queries is now correct

## [0.1.2]

First "working" version

### Fixed

- Fixed: Plugin files in the zip file are now executable
