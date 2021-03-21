# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

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
