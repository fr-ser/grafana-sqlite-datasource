# Release Notes

This is an official release of the Grafana SQLite plugin code. For some releases it takes a while
to also appear in the Grafana repository, which is used to populate the plugins for the Grafana
website as well as their grafana-cli. Some releases (especially pre-releases) from Github never
appear in the Grafana repository.

To install this specific release use one of the following methods.

More details can also be found in the installation section of the [Readme](README.md).

## Installation

### Using the grafana-cli

1. Run this command:

   ```sh
   grafana-cli --pluginUrl https://github.com/fr-ser/grafana-sqlite-datasource/releases/download/v$VERSION/frser-sqlite-datasource-$VERSION.zip plugins install frser-sqlite-datasource
   ```

2. See the installation instructions in the [Readme](README.md).

### Manual Installation

1. Download the [zip file](https://github.com/fr-ser/grafana-sqlite-datasource/releases/download/v$VERSION/frser-sqlite-datasource-$VERSION.zip) below
2. Extract the zip file into the data/plugins subdirectory for Grafana: `unzip <the_download_zip_file> -d <plugin_dir>/`

   Finding the plugin directory can sometimes be a challenge as this is platform and settings
   dependent. A common location for this on Linux devices is `/var/lib/grafana/plugins/`
3. See the installation instructions in the [Readme](README.md).

## Changelog

For the full changelog see [here](CHANGELOG.md).
