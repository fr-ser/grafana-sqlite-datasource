# Examples

This file contains some typical examples for queries to get started.

## Interpolating filter variables

To set up variables by which users can easily filter dashboards, go into the dashboard settings, and click the "variables" option.
You can populate this with queries, for example:

```sql
SELECT name FROM students ORDER BY name ASC;
```

Supposing the name of the variable created with the above query was `students`, you can then easily interpolate this value with an `IN` clause in your dashboard's query:

```sql
SELECT unix_timestamp, present FROM class_arrival_times
WHERE student IN (${students:singlequote})
```

To learn more about advanced variable interpolation to facilitate your queries, see
[here](https://grafana.com/docs/grafana/latest/variables/advanced-variable-format-options/).

## Filter by time

The key here is to use the Grafana variables `__to` and `__from` and format them correctly. By
default those variables represent milliseconds instead of seconds in Unix time.

```sql
SELECT avg(value), max(value), min(value) FROM sine_wave
WHERE time >= $__from / 1000 and time < $__to / 1000
```

```sql
SELECT * FROM products WHERE created_at BETWEEN ${__from:date:seconds} AND ${__to:date:seconds};
```

## Creating a time series

In theory a table can be easily displayed as a time series. There are some considerations to make
it more robust though.

### Dealing with "duplicate readings"

Sometimes there are duplicate readings in a time frame. When ignoring this the chart can seem
skewed (visually and in calculations) as the time frame with more readings is wider (visually) and
overrepresented (in calculations using the average).

The solution is to aggregate time frames. Even if only one reading per frame is expected this
solves the issue of duplicates.

```sql
SELECT (time / 500) * 500  as window, avg(value)
FROM sine_wave
GROUP BY 1
ORDER BY 1 ASC
```

```sql
SELECT $__unixEpochGroupSeconds(time, 500) as window, avg(value)
FROM sine_wave
GROUP BY 1
ORDER BY 1 ASC
```

### Dealing with "missing readings"

Sometimes there are no readings for a time frame but we still want to display this time frame.
In such cases we need to "fill the gaps" of missing readings. They can be filled with `NULL` values.

```sql
SELECT $__unixEpochGroupSeconds(time, 500, NULL) as window, avg(value)
FROM sine_wave
GROUP BY 1
ORDER BY 1 ASC
```

## Convert time

The important part is to generate timestamps that can be recognized by the plugin
[see information in the Readme](https://github.com/fr-ser/grafana-sqlite-datasource#support-for-time-formatted-columns).

The column that represents a timestamp also needs to be selected as a time column in the UI.

### Convert dates to timestamps

```sql
WITH converted AS (
   -- a row looks like this (value, date): 1.45, '2020-12-12'
   SELECT value,  date || 'T00:00:00Z' AS datetime FROM raw_table
)
SELECT datetime, value FROM converted ORDER BY datetime ASC
```
