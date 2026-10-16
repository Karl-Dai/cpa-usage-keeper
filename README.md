<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/keeper-logo-dark.svg" />
    <source media="(prefers-color-scheme: light)" srcset="./assets/keeper-logo-light.svg" />
    <img src="./assets/keeper-logo-light.svg" alt="Keeper" width="560" />
  </picture>
</p>

<p align="center">
  <a href="./README.md"><strong>English</strong></a> ｜ <a href="./README.zh.md">简体中文</a>
</p>

<h1 align="center">CPA Usage Keeper</h1>

<p align="center">Every flow leaves a trace.</p>

<p align="center">
  <a href="https://github.com/Willxup/cpa-usage-keeper/releases/latest"><img src="https://img.shields.io/github/v/release/Willxup/cpa-usage-keeper?style=flat-square" alt="Latest release" /></a>
  <a href="https://github.com/Willxup/cpa-usage-keeper/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/Willxup/cpa-usage-keeper/ci.yml?branch=main&amp;style=flat-square&amp;label=CI" alt="CI status" /></a>
  <a href="https://github.com/Willxup/cpa-usage-keeper/pkgs/container/cpa-usage-keeper"><img src="https://img.shields.io/badge/Docker-GHCR-2496ED?style=flat-square&amp;logo=docker&amp;logoColor=white" alt="Docker image on GHCR" /></a>
  <a href="https://github.com/Willxup/homebrew-cpa-usage-keeper"><img src="https://img.shields.io/badge/Homebrew-supported-FBB040?style=flat-square&amp;logo=homebrew&amp;logoColor=black" alt="Homebrew supported" /></a>
  <a href="https://github.com/Willxup/cpa-usage-keeper/releases/latest"><img src="https://img.shields.io/badge/Linux-FCC624?style=flat-square&amp;logo=linux&amp;logoColor=black" alt="Linux supported" /></a>
  <a href="https://github.com/Willxup/cpa-usage-keeper/releases/latest"><img src="https://img.shields.io/badge/macOS-A2AAAD?style=flat-square&amp;logo=apple&amp;logoColor=black" alt="macOS supported" /></a>
  <a href="https://github.com/Willxup/cpa-usage-keeper/releases/latest"><img src="https://img.shields.io/badge/Windows-0078D4?style=flat-square&amp;logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCI+PHBhdGggZmlsbD0iI2ZmZiIgZD0iTTIgMy41IDExIDJ2OUgyem0xMC0xLjdMMjIgLjNWMTFIMTJ6TTIgMTJoOXY5TDIgMTkuNXptMTAgMGgxMHYxMC43bC0xMC0xLjV6Ii8+PC9zdmc%2B" alt="Windows supported" /></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/Willxup/cpa-usage-keeper?style=flat-square" alt="MIT License" /></a>
</p>

Keep usage history for [CLIProxyAPI (CPA)](https://github.com/router-for-me/CLIProxyAPI), understand model costs, and track request performance and credential quotas. CPA Usage Keeper brings overviews, realtime diagnostics, request details, and quota history into one independently deployed dashboard.

**[Add Keeper to an existing CPA](#keeper-only) · [Deploy CPA + Keeper](#cpa--keeper) · [Configuration](#configuration)**

## Screenshots

<table>
  <tr>
    <td width="50%" align="center">
      <strong>Overview</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/overview-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/overview-white.png" />
        <img src="./assets/screenshots/overview-white.png" alt="CPA Usage Keeper Overview" width="100%" />
      </picture>
    </td>
    <td width="50%" align="center">
      <strong>Realtime Diagnostics</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/realtime-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/realtime-white.png" />
        <img src="./assets/screenshots/realtime-white.png" alt="CPA Usage Keeper Realtime Diagnostics" width="100%" />
      </picture>
    </td>
  </tr>
  <tr>
    <td width="50%" align="center">
      <strong>Usage Analysis</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/analysis-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/analysis-white.png" />
        <img src="./assets/screenshots/analysis-white.png" alt="CPA Usage Keeper Usage Analysis" width="100%" />
      </picture>
    </td>
    <td width="50%" align="center">
      <strong>Request Events</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/request-events-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/request-events-white.png" />
        <img src="./assets/screenshots/request-events-white.png" alt="CPA Usage Keeper Request Events" width="100%" />
      </picture>
    </td>
  </tr>
  <tr>
    <td width="50%" align="center">
      <strong>Credentials and Quotas</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/credentials-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/credentials-white.png" />
        <img src="./assets/screenshots/credentials-white.png" alt="CPA Usage Keeper Credentials and Quotas" width="100%" />
      </picture>
    </td>
    <td width="50%" align="center">
      <strong>Quota History</strong>
      <picture>
        <source media="(prefers-color-scheme: dark)" srcset="./assets/screenshots/quota-history-dark.png" />
        <source media="(prefers-color-scheme: light)" srcset="./assets/screenshots/quota-history-white.png" />
        <img src="./assets/screenshots/quota-history-white.png" alt="CPA Usage Keeper Quota History" width="100%" />
      </picture>
    </td>
  </tr>
</table>

## Features

- **Keep your history**: persist CPA usage in SQLite with scheduled backups.
- **Understand costs**: analyze usage, caching, and estimated costs by model, API Key, and provider.
- **Investigate requests**: inspect and export request details, and explore success rates, time to first token (TTFT), and total latency.
- **Track quotas**: monitor credential health and remaining quotas, refresh quotas, edit priorities, and explore Codex quota history.
- **Share scoped access**: provide a read-only usage view for an individual CPA API Key.

CPA plugin embedding in CPAMC is also available. Deploy with Docker Compose, Homebrew, or binaries; login protection is enabled by default.

## Sponsors and Special Thanks

- Thanks to [CLIProxyAPI (CPA)](https://github.com/router-for-me/CLIProxyAPI) for providing the upstream CPA foundation and data source this project builds on.
- Thanks to [@YouShouldBetOnMe](https://github.com/YouShouldBetOnMe) for supporting CPA Usage Keeper.
- Thanks to the CPA discussion group for their discussions and feedback.

## Quick Start

> Keeper supports CPA v8.0.12+.

> When multiple usage collectors share one CPA instance, ensure they all use subscription mode; otherwise, collection may stop or become incomplete.

Docker Compose is the recommended deployment method. Use the full stack when deploying CPA and Keeper together, or the Keeper-only stack when CPA already exists.

Docker deployment requires Docker and Docker Compose; you do not need Go, Node.js, or a source build. For an existing CPA, have its address and management key ready, and choose a Keeper login password.

| Setup | Recommended path | Architectures |
| --- | --- | --- |
| New CPA + Keeper deployment | [Docker Compose: CPA + Keeper](#cpa--keeper) | `linux/amd64`, `linux/arm64` |
| Existing CPA deployment | [Docker Compose: Keeper only](#keeper-only) | `linux/amd64`, `linux/arm64` |
| Existing CPA, Docker CLI preferred | [Docker](#docker-cpa-already-runs-on-the-host) | `linux/amd64`, `linux/arm64` |
| macOS | [Homebrew](#macos-homebrew) | `amd64`, `arm64` |
| Linux without containers | [Linux binary](#linux-binary) | `amd64`, `arm64` |
| Windows | [Windows binary](#windows-binary) | `amd64`, `arm64` |

Login protection is enabled by default. Configure `LOGIN_PASSWORD` before starting Keeper, or explicitly set `AUTH_ENABLED=false` only when access is reliably isolated by the deployment environment.

## Benchmark

Production-style `linux/amd64` capacity measurements for sustained ingestion, Dashboard latency, CPU utilization, and Keeper cgroup peak memory are available in the [Capacity Benchmark Report](./internal/benchmark/reports/capacity/README.md).

<details>
<summary>Developer reference: project structure, local setup, and tests</summary>

## Project Structure

```text
cmd/server/              Application entry point
internal/api/            HTTP routes and handlers
internal/app/            Application wiring and startup
internal/auth/           Sessions and access control
internal/poller/         CPA usage and metadata synchronization
internal/repository/     SQLite persistence and aggregations
internal/service/        Usage, pricing, and identity services
internal/quota/          Provider quota refresh and inspection
internal/benchmark/      Capacity suite, reports, manifests, and legacy microbenchmarks
deploy/                  Deployment templates
web/                     React + TypeScript frontend
```

## Local Development

### Prerequisites

- Go 1.26+
- Node.js 24+
- npm
- A running [CLIProxyAPI (CPA)](https://github.com/router-for-me/CLIProxyAPI) instance

### Run Locally

1. Copy `.env.example` to `.env`, then set at least `CPA_BASE_URL`, `CPA_MANAGEMENT_KEY`, and a private `LOGIN_PASSWORD`.

```bash
cp .env.example .env
vim .env
```

2. Start the backend.

```bash
go run ./cmd/server/main.go
```

3. In another terminal, install frontend dependencies and start the development server.

```bash
npm --prefix ./web ci
npm --prefix ./web run dev -- --host 127.0.0.1
```

Open `http://127.0.0.1:5173`. The frontend proxies `/api` to `http://127.0.0.1:8080`; override it with `VITE_API_PROXY_TARGET` when the backend uses another port.

### Tests

Run the full verification baseline:

```bash
make verify
```

Or run checks individually:

```bash
go test ./cmd/... ./internal/...
npm --prefix ./web run test
npm --prefix ./web run lint
npm --prefix ./web run typecheck
npm --prefix ./web run build
```

</details>

## Deployment

After startup, open `http://your-server-address:8080` (`http://127.0.0.1:8080` for a local deployment) and sign in with your Keeper login password. Adjust the URL if you configure a different port, HTTPS, or a base path.

### Docker Compose (Recommended)

Docker Compose is recommended for both a complete CPA + Keeper stack and a Keeper-only deployment.

#### CPA + Keeper

**1. Prepare the configuration files**

First, download the configuration example from the [official CPA repository](https://github.com/router-for-me/CLIProxyAPI) into `./cpa/config.yaml` in your deployment directory:

```bash
mkdir -p cpa/auths cpa/logs keeper
curl -fL https://raw.githubusercontent.com/router-for-me/CLIProxyAPI/main/config.example.yaml \
  -o cpa/config.yaml
```

In the same directory, download the [CPA + Keeper Compose template](./deploy/docker-compose.full.example.yml):

```bash
curl -fL https://raw.githubusercontent.com/Willxup/cpa-usage-keeper/main/deploy/docker-compose.full.example.yml \
  -o docker-compose.yml
```

Use this download command for a new deployment only; edit an existing configuration directly to avoid overwriting it.

**2. Configure CPA and Keeper**

Edit the existing settings in `cpa/config.yaml` as listed below, preserving YAML nesting and indentation. Dots in the table denote nested paths, not literal YAML keys to add:

| Setting | What to configure |
| --- | --- |
| `management.allow-remote` | Set to `true` so Keeper can access CPA management features from another container. |
| `management.secret-key` | Set a private management key and use the same original key for `CPA_MANAGEMENT_KEY` in `docker-compose.yml`. |
| `observability.usage.usage-statistics-enabled` | Set to `true` to enable usage statistics. |
| `access.api-keys` | Replace the example keys with your own client API keys; these are separate from the CPA management key and Keeper login password. |

Keep `server.host` as an empty string, `server.port` as `8317`, and `oauth.auth-dir` as `"~/.cli-proxy-api"` to match the container network and volume mounts in the template.

Edit the downloaded `docker-compose.yml` and fill in two values under Keeper's `environment`:

- `CPA_MANAGEMENT_KEY`: use the same original management key configured in CPA.
- `LOGIN_PASSWORD`: replace the empty string `""` with your own Keeper login password.

Alternatively, put the settings in `./keeper/.env`; the template loads it through an optional `env_file` entry and accepts a missing file. Values in `environment` take precedence, so remove the corresponding entries, including empty placeholders, to use values from the file.

**3. Start and open Keeper**

```bash
docker compose up -d
```

Open `http://your-server-address:8080` and sign in with your Keeper password. Run `docker compose down` to stop the stack.

For a new deployment, add model credentials in CPA and make a model request to generate usage records.

CPA data is stored under `./cpa`, and Keeper data is stored under `./keeper`.

#### Keeper Only

**1. Prepare the configuration files**

When CPA is already deployed, download the Keeper-only Compose template and environment configuration into a new deployment directory:

```bash
curl -fL https://raw.githubusercontent.com/Willxup/cpa-usage-keeper/main/deploy/docker-compose.example.yml \
  -o docker-compose.yml
curl -fL https://raw.githubusercontent.com/Willxup/cpa-usage-keeper/main/.env.example -o .env
vim .env
```

For an existing deployment, edit the existing files directly to avoid overwriting your configuration.

**2. Set the connection details and login password**

For CPA running on the Docker host, start with:

```env
CPA_BASE_URL=http://host.docker.internal:8317
CPA_MANAGEMENT_KEY=replace-with-your-management-key
AUTH_ENABLED=true
LOGIN_PASSWORD=
```

Set a private `LOGIN_PASSWORD` before starting the container.

Whether CPA runs on the Docker host or another machine, it must allow remote management and listen on an address reachable from the Keeper container; a listener bound only to `127.0.0.1` is not reachable. Set `management.allow-remote: true` and check `server.host`.

For other network layouts, set `CPA_BASE_URL` to a CPA address reachable from the container. Set `REDIS_QUEUE_ADDR` only when the Redis/RESP address differs from the automatically derived address.

**3. Start and open Keeper**

```bash
docker compose up -d
```

Open `http://your-server-address:8080` and sign in with your Keeper password. Run `docker compose down` to stop the stack.

Keeper data is stored under `./data` by the provided template.

#### Logs and Updates

For either Compose setup, view Keeper logs from the deployment directory:

```bash
docker compose logs --tail=100 -f cpa-usage-keeper
```

For the first V1 → V2 upgrade, follow [Upgrade Safety and Recovery](#upgrade-safety-and-recovery) first to create an independent snapshot.

For routine updates, keep your configuration and data directories and run:

```bash
docker compose pull cpa-usage-keeper
docker compose up -d cpa-usage-keeper
```

#### No Data After Startup?

- Check that CPA usage statistics are enabled: `observability.usage.usage-statistics-enabled` must be `true`.
- Ensure Keeper can reach `CPA_BASE_URL` and that `CPA_MANAGEMENT_KEY` matches the CPA management key.
- Check that CPA has received new model requests. If data is still missing, use the log command above to look for connection or authentication errors.

### Docker (CPA Already Runs On The Host)

Use the same `.env` values as the Keeper-only Compose setup when you prefer `docker run`:

```bash
docker run -d \
  --name cpa-usage-keeper \
  --add-host=host.docker.internal:host-gateway \
  -p 8080:8080 \
  -v "$(pwd)/keeper:/data" \
  --env-file .env \
  ghcr.io/willxup/cpa-usage-keeper:latest
```

### macOS Homebrew

Homebrew is the recommended macOS installation method:

```bash
brew tap Willxup/cpa-usage-keeper
brew install cpa-usage-keeper
```

Set `CPA_BASE_URL`, `CPA_MANAGEMENT_KEY`, and a private `LOGIN_PASSWORD`, then start the service:

```bash
vim "$(brew --prefix)/etc/cpa-usage-keeper.env"
brew services start cpa-usage-keeper
```

Upgrade and service commands:

```bash
brew services list
brew services restart cpa-usage-keeper
brew update
brew upgrade cpa-usage-keeper
```

Data is stored under `$(brew --prefix)/var/cpa-usage-keeper`; logs are written under `$(brew --prefix)/var/log/`.

### Linux Binary

Download the `linux_amd64` or `linux_arm64` archive from [Releases](https://github.com/Willxup/cpa-usage-keeper/releases/latest), then extract and run it:

```bash
mkdir -p cpa-usage-keeper
tar -xzf ./cpa-usage-keeper_*_linux_*.tar.gz -C cpa-usage-keeper --strip-components=1
cd cpa-usage-keeper
cp .env.example .env
vim .env
./cpa-usage-keeper
```

#### systemd

The Linux package includes a service template. Run these commands from the extracted package directory:

```bash
sudo cp cpa-usage-keeper.service /etc/systemd/system/cpa-usage-keeper.service
sudo sed -i "s|__CPA_USAGE_KEEPER_DIR__|$(pwd)|g" /etc/systemd/system/cpa-usage-keeper.service
sudo systemctl daemon-reload
sudo systemctl enable --now cpa-usage-keeper
```

```bash
sudo systemctl status cpa-usage-keeper
sudo journalctl -u cpa-usage-keeper -f
sudo systemctl restart cpa-usage-keeper
```

### Command-Line Options

The binary supports optional startup flags:

```bash
cpa-usage-keeper --env /path/to/keeper.env # Use a specific environment file.
cpa-usage-keeper --host 127.0.0.1 # Override APP_HOST for this process.
cpa-usage-keeper -v               # Print the build version and exit; --version is also supported.
```

### Windows Binary

Download the `windows_amd64` or `windows_arm64` ZIP package from [Releases](https://github.com/Willxup/cpa-usage-keeper/releases/latest) and extract it. In PowerShell, open the extracted package directory and run:

```powershell
Copy-Item .env.example .env
notepad .env
.\cpa-usage-keeper.exe
```

Set `CPA_BASE_URL`, `CPA_MANAGEMENT_KEY`, and a private `LOGIN_PASSWORD` before starting. Authentication is enabled by default; set `AUTH_ENABLED=false` explicitly only for an isolated deployment.

### Upgrade Safety and Recovery

Before the first V1 → V2 upgrade, stop Keeper and take an independent snapshot of the entire `WORK_DIR`, including `app.db` and any SQLite `-wal`/`-shm` files. Keep the snapshot outside `WORK_DIR/backups` and retain the V1 binary/image and configuration. The migration's automatic M1 backup uses the normal backup directory and is subject to daily backup retention; it is not a permanent rollback snapshot.

Set `keeper_data_dir` to the host directory mounted at `/data`: the Keeper-only template uses `./data`, while the CPA + Keeper template uses `./keeper`. Adjust the value below to match your deployment and choose a snapshot directory that does not already exist. The chained commands stop if any step fails:

```bash
keeper_data_dir=./data # CPA + Keeper: use ./keeper instead
test -f "$keeper_data_dir/app.db" &&
test ! -e ./keeper-v1-snapshot &&
docker compose stop cpa-usage-keeper &&
cp -a "$keeper_data_dir" ./keeper-v1-snapshot &&
docker compose pull cpa-usage-keeper &&
docker compose up -d --no-deps cpa-usage-keeper
```

After Keeper is ready, it checks CPA once in the background. An explicitly identified release older than v8.0.12 produces an upgrade warning in the server log. Supported or unknown versions and unavailable CPA do not produce a version warning. This check never blocks startup, local history access, or data ingestion; remote features still depend on the actual CPA API compatibility.

If startup fails, inspect the logs and `/api/v1/startup/status` (with `APP_BASE_PATH` when configured). Correcting the cause and restarting resumes the migration. When moving `WORK_DIR`, preserve `backups/YYYY-MM-DD/database_TIMESTAMP.db`: if the original M1 path is missing and its SHA-256 was recorded, Keeper can verify that exact date/file name under the new `WORK_DIR/backups` and resume. There is no `BACKUP_DIR` environment variable.

**Offline restore is only for an incomplete legacy migration** (`pricing_migration_state.data_complete = 0`). A startup failure after migration has completed is not eligible: normal inbox cleanup may already have removed messages needed for this recovery method. For a completed migration, use an independently retained snapshot and its matching application version instead. A snapshot rollback restores that point in time and does not retain data received after the snapshot.

1. Stop Keeper and every process writing its database. Preserve the failed `WORK_DIR`, including existing SQLite sidecars, and the M1 backup before restoring. In the commands below, `keeper_data_dir` must point to that failed data directory (set it again if using a new shell).
2. Locate the M1 `backup_path` recorded in the failed database. With the optional SQLite CLI, inspect it read-only:

   ```bash
   sqlite3 -readonly "$keeper_data_dir/app.db" \
     'SELECT data_complete, backup_path FROM pricing_migration_state WHERE id = 1;'
   ```

3. Run the V2 `pricing-restore` tool, supplied beside the server in binary packages and installed by Homebrew. Set `-backup` to the original M1 file or an unchanged copy, `-fault` to the stopped failed database, and `-out` to a new file that does not exist:

   ```bash
   mkdir -p ./recovery
   ./pricing-restore \
     -backup /path/to/M1-copy.db \
     -fault "$keeper_data_dir/app.db" \
     -out ./recovery/restored.db
   ```

   On Windows use `.\pricing-restore.exe`; with Homebrew use `pricing-restore`. The tool verifies the recorded SHA-256, copies the M1 snapshot, and replays later inbox rows by their original IDs as pending. It publishes the new database only after validation and leaves both inputs unchanged. A backup with a nonempty WAL/journal is rejected: use the original standalone M1 file. For older unfinished states without a recorded fingerprint, the original recorded path must still be readable so the selected copy can be compared; if that evidence is lost, the tool refuses to guess.

For Docker, keep the original service stopped. Set `keeper_data_dir` to the same actual data directory as above (`./data` for Keeper-only, `./keeper` for CPA + Keeper). This example mounts it read-only and writes only to a separate recovery directory. Replace the image tag with the V2 build used for the upgrade and the example M1 path with the recorded date/file name:

```bash
mkdir -p ./recovery
docker run --rm --entrypoint /app/pricing-restore \
  --mount type=bind,src="$(cd -- "$keeper_data_dir" && pwd)",dst=/data,readonly \
  --mount type=bind,src="$(pwd)/recovery",dst=/restore \
  ghcr.io/willxup/cpa-usage-keeper:YOUR_V2_TAG \
  -backup /data/backups/YYYY-MM-DD/database_TIMESTAMP.db \
  -fault /data/app.db \
  -out /restore/restored.db
```

To use a relocated backup outside that mount, mount the copy read-only as well and pass its container path to `-backup`. The stored fingerprint identifies the file independently of its name or location. An older state without a fingerprint still needs its original recorded path available inside the container.

After checking the tool's inbox boundaries and replay count, place the result in a new data directory and run the matching V1 version against that directory:

```bash
mkdir ./restored-data
mv ./recovery/restored.db ./restored-data/app.db
```

For Compose, map `./restored-data:/data` and use the retained V1 image; for a binary/Homebrew service, set `WORK_DIR` to the new directory. Keep the failed directory intact. Do not place the restored main file beside the failed database's old WAL/SHM files.

### Health Checks

`GET /healthz` is a liveness check: it returns `200` while Keeper can serve HTTP, including during migration or after startup fails. Use `GET /readyz` for deployment readiness: it returns `503` while opening, migrating, or failed, and `200` only after the business service is ready. Both routes use `APP_BASE_PATH`, for example `/cpa/readyz`; readiness responses are not cached.

The Docker image checks `/readyz`. A long migration can mark the container `unhealthy`; Docker itself does not restart a container because of that health status. In an orchestrator, use `/readyz` for readiness and `/healthz` for liveness. Do not restart on failed readiness or impose a short startup-probe deadline on database migrations. A failed startup keeps its status page available; inspect the logs, correct the cause, and restart Keeper to retry.

## Configuration

Copy the example config:

```bash
cp .env.example .env
```

For a first deployment, set the CPA address, CPA management key, and Keeper login password. Most other settings can keep their defaults. Read the relevant sections when you need a domain, HTTPS, or a base path.

### Minimum Required

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `CPA_BASE_URL` | Yes | - | URL used by the Keeper server to call CPA. In Docker Compose this is usually `http://cli-proxy-api:8317`, and it can be a private address or container service name |
| `CPA_MANAGEMENT_KEY` | Yes | - | CPA management key used to read CPA management APIs |

### Web Access And Reverse Proxy

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `APP_HOST` | No | all interfaces | Keeper HTTP listen host; native deployments can set `127.0.0.1` for local-only access |
| `APP_PORT` | No | `8080` | Keeper HTTP listen port |
| `APP_BASE_PATH` | No | root path | Keeper subpath prefix, such as `/keeper`; empty means `/` |
| `CPA_PUBLIC_URL` | No | current browser origin root | Public CPA URL for the "Back to CPA" link and CPAMC frame trust |
| `TRUSTED_PROXY_CIDRS` | No | local loopback only | Additional reverse-proxy CIDRs allowed to provide `X-Forwarded-For`, separated by commas |

- The `--host` startup flag overrides `APP_HOST`. When neither is set, Keeper preserves its existing behavior and listens on all available network interfaces.
- For Docker/Compose, keep `APP_HOST` empty. To restrict access to the Docker host, publish the port as `127.0.0.1:8080:8080`.
- `APP_BASE_PATH` must be empty or start with `/`; `/cpa/` is normalized to `/cpa`.
- `CPA_BASE_URL` is the server-side CPA address and may use a private host or Docker service name.
- `CPA_PUBLIC_URL` controls browser navigation and cross-origin CPAMC frame trust. Leave it empty for same-origin `/management.html`, or set an explicit public CPA URL when domains, ports, or paths differ.
- Keeper trusts `X-Forwarded-For` only from local loopback and `TRUSTED_PROXY_CIDRS`. Direct clients cannot change their login-rate-limit source with this header. Configure only the exact proxy address or network; universal CIDRs are rejected.

For cross-origin CPAMC embedding, `CPA_PUBLIC_URL` must be a complete `http://` or `https://` URL with a host. Relative paths affect navigation only.

### Login Protection

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `AUTH_ENABLED` | No | `true` | Enable login protection |
| `LOGIN_PASSWORD` | When auth is enabled | - | Login password |
| `CPA_REQUEST_LOG_ACCESS_ENABLED` | No | `false` | Allow administrators to view and download CPA request logs through Keeper; corresponding logs must exist in CPA and may contain request or response data |
| `AUTH_SESSION_TTL` | No | `168h` | Login session lifetime |

### Timezone And Request Behavior

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `TZ` | No | `Asia/Shanghai` | Timezone used for statistics and display; Today, daily totals, page timestamps, log timestamps, and daily cleanup are calculated in this timezone |
| `REQUEST_TIMEOUT` | No | `30s` | Timeout for CPA HTTP requests and Redis queue operations |
| `TLS_SKIP_VERIFY` | No | `false` | Skip TLS certificate verification for CPA HTTPS and Redis queue TLS; enable only with self-signed certificates |

### Auth Files Quota Refresh

Scheduled Auth Files quota refresh is configured from the gear button in the Auth Files inspection dialog. The setting is stored in the local SQLite database and does not require the page to stay open.

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `QUOTA_REFRESH_WORKER_LIMIT` | No | `10` | Maximum Auth Files quota refresh concurrency for manual and scheduled refresh, capped at `100` |
| `QUOTA_UPSTREAM_RESPONSES_ENABLED` | No | `false` | Cache each credential's latest raw upstream quota responses and return them through quota task/cache APIs for browser Network debugging; responses may contain account data |

### Redis Queue Advanced Settings

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `REDIS_QUEUE_ADDR` | No | Derived from `CPA_BASE_URL` | An explicit value takes precedence. When empty, use the hostname and explicit port from `CPA_BASE_URL`, or `8317` if no port is specified. Set `host:port` for a separate address |
| `REDIS_QUEUE_TLS` | No | `false` | Use TLS for Redis queue connection; set `true` when `REDIS_QUEUE_ADDR` is explicit and requires TLS |
| `REDIS_QUEUE_BATCH_SIZE` | No | `10000` | Maximum queue records per pull |
| `REDIS_QUEUE_IDLE_INTERVAL` | No | `1s` | Empty queue check interval |

### Storage, Logs, And Backups

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `WORK_DIR` | No | `./data` | Application work directory; database, logs, and backups default to `app.db`, `logs/`, and `backups/` under it |
| `LOG_LEVEL` | No | `info` | Log level |
| `LOG_FILE_ENABLED` | No | `true` | Write persistent log files |
| `LOG_RETENTION_DAYS` | No | `7` | Combined-log history days, plus the current day; `0` disables cleanup. Error-only logs keep 30 history days plus the current day |
| `BACKUP_ENABLED` | No | `true` | Enable SQLite database backups |
| `BACKUP_INTERVAL` | No | `24h` | Database backup interval |
| `BACKUP_RETENTION_DAYS` | No | `7` | Backup retention days |
| `USAGE_RAW_RETENTION_DAYS` | No | `0` | Total raw-request retention in days; only archives are deleted. `0` retains archives forever; `>=90` enables cleanup. Negative values and `1–89` fall back to `0` with a warning without preventing startup |

Keeper archives raw request records older than 90 local calendar days daily at 04:30 in the configured timezone. Archives are retained permanently by default. With `USAGE_RAW_RETENTION_DAYS>=90`, daily maintenance deletes expired archived records in batches after archiving, using the request timestamp and local calendar days in the configured timezone. The 90-day hot-table window remains unchanged. Normal pages do not directly query archived records; historical aggregates and backups keep their own retention policies. Deleted raw records cannot support later historical recalculation. Deletion may not immediately shrink the SQLite file: freed space can be reused and is reclaimed by the existing conditional compaction policy.

When file logging is enabled, `cpa-usage-keeper-YYYY-MM-DD.log` contains all emitted levels. Error, fatal, and panic entries are also copied to `cpa-usage-keeper-error-YYYY-MM-DD.log`, which keeps the previous 30 local calendar dates plus the current date.

### Built-In HTTPS

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `TLS_ENABLED` | No | `false` | Let Keeper serve HTTPS/TLS directly |
| `TLS_CERT_FILE` | Required when TLS is enabled | - | HTTPS certificate file path |
| `TLS_KEY_FILE` | Required when TLS is enabled | - | HTTPS private key file path |

Usually, HTTPS should be terminated at nginx, Caddy, or another reverse proxy. Set `TLS_ENABLED=true` only when the Keeper process must serve HTTPS directly, and provide `TLS_CERT_FILE` and `TLS_KEY_FILE`; relative paths are resolved against the `.env` file directory.

Security and data notes:

- Browser APIs redact key-like fields, but the SQLite database and its unencrypted backups contain original data.
- Authentication is enabled by default. If it is explicitly disabled, restrict Keeper access at the deployment boundary; terminate public HTTPS at a reverse proxy.
- Login session hashes persist in SQLite until logout or `AUTH_SESSION_TTL` expiry.
- CPAMC uses a separate embed session: an `HttpOnly` cookie when available, or a per-tab header token in browser session storage as a fallback.
- Same-origin embedding works by default. For cross-origin embedding, set `CPA_PUBLIC_URL` to the public CPA/CPAMC origin used for `frame-ancestors`.
- Redis inbox messages are retained through the current day after success or for 7 days after failure.

## Nginx Reverse Proxy

When serving under `/cpa`, set `APP_BASE_PATH=/cpa` and keep the prefix in your reverse proxy:

```nginx
location /cpa/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

The loopback Nginx configuration above works without additional Keeper settings. If the reverse proxy reaches Keeper from a container or another host, or a CDN such as Cloudflare sits in front of the reverse proxy, add the exact proxy networks, for example `TRUSTED_PROXY_CIDRS=172.18.0.0/16`.

When CPA and Keeper share a browser origin, `CPA_PUBLIC_URL` can be omitted and "Back to CPA" uses `/management.html`. For another domain, port, or path, set the public CPA URL:

```env
CPA_PUBLIC_URL=https://cpa.example.com
```

## Benchmark

Production-style `linux/amd64` capacity measurements for sustained ingestion, Dashboard latency, CPU utilization, and Keeper cgroup peak memory are available in the [Capacity Benchmark Report](./internal/benchmark/REPORT.md).

## License

This project is open source under the [MIT License](./LICENSE).
