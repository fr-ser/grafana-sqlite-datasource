# Common Problems / FAQ

The following section describes common issues encountered while using the plugin.

## I have a "file not found" error for my database

The first choice should be to make sure, that the path is correct. It is also good practice to
use absolute paths (e.g. `/app/state/data.db`) instead of relative paths (`state/data.db`).

In case the path is correct but the database is in the `/var` directory on a linux system there
might also be a systemd issue. This is typically observed with Grafana versions starting with
v8.2.0. When Grafana is run via systemd (the typical default installation on Linux systems) the
`/var` directory is not available to Grafana (and therefore also not to the plugin).

In order to change this behavior you need to do the following:

```txt
# edit (override) the grafana systemd configuration
systemctl edit grafana-server

# add the following lines
[Service]
PrivateTmp=false

# reload the systemd config and restart the app
systemctl daemon-reload
systemctl restart grafana-server
```

## I have a "permission denied" error for my database

Make sure, that you have access to the file and all the folders in the path of the file.
Read access is enough for the plugin.

In case the permissions are correct but database is in the `/home` directory on a linux system
there might also be a systemd issue. This is typically observed with Grafana versions starting with
v8.2.0. When Grafana is run via systemd (the typical default installation on Linux systems) the
`/home` directory is not available to Grafana (and therefore also not to the plugin).

In order to change this behavior you need to do the following:

```txt
# edit (override) the grafana systemd configuration
systemctl edit grafana-server

# add the following lines
[Service]
ProtectHome=false

# reload the systemd config and restart the app
systemctl daemon-reload
systemctl restart grafana-server
```

## The legend of my time series appears twice / is doubled

Sometimes (especially when displaying multiple lines in a time series chart) the legend (the information below the chart) can show the name of the column twice.
The legend can read "value value" or "temperature temperature".

This can be controlled through the field display name configuration.
There a hardcoded value can be set but the value can also be based on the "labels" of the search.
Some more information about setting the display name via labels can be found in the [Grafana documentation](https://grafana.com/docs/grafana/latest/panels/configure-standard-options/#display-name).
