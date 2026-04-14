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

## I want to open a read only database and get errors

If you get an error like `attempt to write a readonly database`, the most common cause is that
your database is in **WAL (Write-Ahead Logging) mode** and Grafana does not have write access to
the directory containing the database file.

In WAL mode, SQLite must create or update a shared-memory index file (`<db>-shm`) alongside the
database — even for read-only connections. If the directory is not writable by the Grafana
process, SQLite cannot create this file and the connection fails, regardless of `mode=ro` being
set in the path options.

To check whether your database uses WAL mode:

```sh
sqlite3 /path/to/your.db "PRAGMA journal_mode;"
# returns "wal" if WAL mode is active
```

**Option 1 — switch to journal mode** (simplest, if you control the database):

```sh
sqlite3 /path/to/your.db "PRAGMA journal_mode=DELETE;"
```

**Option 2 — ensure the directory is writable** by the Grafana process so that SQLite can manage
the `-shm` file. Read access on the database file itself is sufficient; only the directory needs
to be writable.

**Option 3 — use `immutable=1`** in the path options if the database will never change while
Grafana is running. This bypasses all WAL locking requirements entirely:

```txt
immutable=1
```

> Warning: with `immutable=1`, SQLite will not see any changes made to the database after it is
> opened. Only use this if the database is truly static.

For more background see the official SQLite documentation on
[read-only WAL databases](https://www.sqlite.org/wal.html#read_only_databases).

## The legend of my time series appears twice / is doubled

Sometimes (especially when displaying multiple lines in a time series chart) the legend (the information below the chart) can show the name of the column twice.
The legend can read "value value" or "temperature temperature".

This can be controlled through the field display name configuration.
There a hardcoded value can be set but the value can also be based on the "labels" of the search.
Some more information about setting the display name via labels can be found in the [Grafana documentation](https://grafana.com/docs/grafana/latest/panels/configure-standard-options/#display-name).

## Can I run the plugin with Grafana Cloud

Currently (2024-03-08) there is little use in running the SQLite plugin with Grafana Cloud.

The problem is that Grafana (and the plugin) run on a separate "cloud instance" and normally SQLite databases are only locally accessible.

It can be useful to run the plugin to connect to an SQLite database on the Grafana cloud instance but that is rarely the goal.

The current ["Grafana Agent"](https://grafana.com/docs/agent/latest/) that is installed on a local machine is only about collecting logs and traces.
No plugins are executed with the agent, which makes it not relevant for this plugin.

## Can I use provisioning with this plugin

Any (backend) plugin supports provisioning; this one included.
The main question is which values to use.
The values can be derived by looking at the configuration of the plugin here:
<https://github.com/fr-ser/grafana-sqlite-datasource/blob/main/src/types.ts>.

An example provisioning file would look like this:

```yaml
apiVersion: 1
datasources:
  - name: sqlite
    type: frser-sqlite-datasource
    access: proxy
    isDefault: true
    editable: true
    jsonData:
      path: /app/data.db
```
