# Plugin Installation

The most up to date (but also most generic) information can always be found here:
[Grafana Website - Plugin Installation](https://grafana.com/docs/grafana/latest/plugins/installation/#install-grafana-plugins)

## Recommended: Installing the Official and Released Plugin on an Existing Grafana With the CLI

Grafana comes with a command line tool that can be used to install plugins.

1. Run this command: `grafana-cli plugins install frser-sqlite-datasource`
2. Restart the Grafana server.
3. To make sure the plugin was installed, check the list of installed data sources. Click the
   Plugins item in the main menu. Both core data sources and installed data sources will appear.

## Latest Version: Installing the newest Plugin Version on an Existing Grafana With the CLI

The grafana-cli can also install plugins from a non-standard URL. This way even plugin versions,
that are not (yet) released to the official Grafana repository can be installed.

1. Run this command:

   ```sh
   # replace the $VERSION part in the URL below with the desired version (e.g. 2.0.2)
   grafana-cli --pluginUrl https://github.com/fr-ser/grafana-sqlite-datasource/releases/download/v$VERSION/frser-sqlite-datasource-$VERSION.zip plugins install frser-sqlite-datasource
   ```

2. See the recommended installation above (from the restart step)

## Manual: Installing the Plugin Manually on an Existing Grafana

In case the grafana-cli does not work for whatever reason plugins can also be installed manually.

1. Get the zip file from [Latest release on Github](https://github.com/fr-ser/grafana-sqlite-datasource/releases/latest)
2. Extract the zip file into the data/plugins subdirectory for Grafana:
   `unzip <the_download_zip_file> -d <plugin_dir>/`

   Finding the plugin directory can sometimes be a challenge as this is platform and settings
   dependent. A common location for this on Linux devices is `/var/lib/grafana/plugins/`
3. See the recommended installation above (from the restart step)
