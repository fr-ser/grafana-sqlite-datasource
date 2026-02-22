# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Grafana backend plugin (`frser-sqlite-datasource`) that allows using an SQLite database as a Grafana data source. It uses `modernc.org/sqlite` (pure Go, no CGO) so the backend compiles without native SQLite libraries.

## Commands

All common tasks are managed via `make`. Run `make help` for a full list.

### Install

```sh
make install            # installs both Go and JS dependencies
make install-go-dependencies   # Go deps + golangci-lint
make install-js-dependencies   # npm install
```

### Build

```sh
make build              # frontend + backend (local arch + docker/linux)
make build-frontend     # webpack production build → dist/
make build-backend-local   # go build for current OS/arch → dist/
make build-backend-docker  # go build for linux (current arch) → dist/
make build-backend-all     # mage: all supported platforms including freebsd, arm
make build-frontend-watch  # webpack --watch (dev mode)
```

Backend binaries are output as `dist/gpx_sqlite-datasource_{os}_{arch}`.

### Run locally

```sh
make start   # docker compose up grafana (credentials: admin / admin123) at http://localhost:3000/
make teardown
```

### Test

```sh
make test-backend    # gotestsum + golangci-lint for ./pkg/...
make test-frontend   # jest (CI mode) + eslint + typecheck
make test-e2e        # build backend for docker + frontend, then selenium tests
make test-e2e-no-build  # selenium tests without rebuilding
make test            # all of the above
```

Run a single Go test:

```sh
gotestsum --format testname -- -count=1 -run TestFunctionName ./pkg/plugin...
```

Run a single frontend test file:

```sh
npx jest --testPathPattern=QueryEditor
```

### Lint / format

```sh
npm run lint:fix     # eslint fix + prettier write
golangci-lint run ./pkg/...
```

## Architecture

The plugin has two halves that communicate via the Grafana plugin protocol (gRPC):

### Go backend (`pkg/`)

Entry point: `pkg/main.go` → `datasource.Manage("frser-sqlite-datasource", plugin.NewDataSource, ...)`

All logic lives in `pkg/plugin/`:

- **`sqlite_datasource.go`** — implements `backend.QueryDataHandler` and `backend.CheckHealthHandler`. `NewDataSource` reads datasource settings (path, pathPrefix, pathOptions, attachLimit) and appends `_pragma=query_only(1)` by default unless `GF_PLUGIN_UNSAFE_DISABLE_QUERY_ONLY_PATH_OPTION=true`.

- **`query.go`** — core query pipeline:

  1. `query()` — entry point per DataQuery: path-block check → variable replacement → macro expansion → `fetchData` → gap filling → build Grafana `data.Frame`
  2. `fetchData()` — opens DB as `{PathPrefix}{Path}?{PathOptions}`, sets `SQLITE_LIMIT_ATTACHED`, runs query, maps column types
  3. Column type mapping: declared DB type (INTEGER/REAL/etc.) → runtime value type inference for `UNKNOWN` columns. Columns listed in `timeColumns` are overridden to `"TIME"`.
  4. For time series queries, long-format frames are converted to wide via `data.LongToWide`, then split into one frame per value field (required for Grafana panel compatibility).

- **`macros.go`** — regex-based SQL macro replacement. Currently supports only `$__unixEpochGroupSeconds(col, interval[, NULL])`, which generates a GROUP BY expression and optionally sets gap-fill config on `queryConfigStruct`.

- **`variables.go`** — replaces `$__from` and `$__to` with millisecond unix timestamps from `dataQuery.TimeRange` (used for alert queries where frontend template variables are unavailable).

- **`gap_filling.go`** — fills missing time-series points with NULL when the `NULL` fill argument is used with `$__unixEpochGroupSeconds`.

- **`block_path.go`** — `IsPathBlocked()` checks a path against a hardcoded security blocklist (`.aws`, `.ssh`, `grafana.db`, etc.) and a user-configurable `GF_PLUGIN_BLOCK_LIST` env var (comma-separated). Called on every query and health check. Controlled by `GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST` and `GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST`.

- **`check_health.go`** — validates the database path and tries to open it.

Security-sensitive env vars (set via `grafana.ini` `[plugin.frser-sqlite-datasource]` section):

- `GF_PLUGIN_UNSAFE_ALLOW_ATTACH_LIMIT_ABOVE_ZERO`
- `GF_PLUGIN_UNSAFE_DISABLE_SECURITY_BLOCKLIST`
- `GF_PLUGIN_UNSAFE_DISABLE_GRAFANA_INTERNAL_BLOCKLIST`
- `GF_PLUGIN_UNSAFE_DISABLE_QUERY_ONLY_PATH_OPTION`
- `GF_PLUGIN_BLOCK_LIST`

### TypeScript/React frontend (`src/`)

- **`module.ts`** — plugin entry; exports `DataSource`, `ConfigEditor`, `QueryEditor`.
- **`types.ts`** — `SQLiteQuery` (`rawQueryText`, `queryText`, `timeColumns`, `queryType`) and `MyDataSourceOptions` (`path`, `pathPrefix`, `pathOptions`, `attachLimit`).
- **`DataSource.ts`** — extends `DataSourceWithBackend`. `applyTemplateVariables` copies `rawQueryText` → `queryText` after template replacement. `metricFindQuery` supports 1-column results (value only) or 2-column results with `__text`/`__value` naming.
- **`QueryEditor.tsx`** — Monaco SQL editor (with legacy textarea fallback), format-as selector (Table / Time series), and time-column tag input.
- **`ConfigEditor.tsx`** — datasource config form: path, pathPrefix (default `file:`), pathOptions (default `_pragma=query_only(1)`), securePathOptions, attachLimit.

### Build system

- **Makefile** — primary task runner; detects OS/arch for binary naming.
- **Magefile.go** — used by `build-backend-all` via `mage`; calls the Grafana SDK's `BuildAll` then adds freebsd/amd64 and linux/arm (ARMv6) builds.
- **Webpack** config is in `.config/webpack/webpack.config.ts` (Grafana plugin scaffold).

### CI (`.github/workflows/branches.yml`)

On every push: backend tests + lint, frontend tests + lint + typecheck, multi-arch backend build, selenium E2E tests. On `main` only: `semantic-release` for automated versioning and tagging.
