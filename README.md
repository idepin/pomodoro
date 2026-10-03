# Pomodoro Backend

REST API backend for a Pomodoro timer application. Built with Go standard library (`net/http`) using a clean architecture layout.

## Requirements

- Go 1.22+ (tested on Go 1.23.4)

## Run

```bash
go run ./cmd/server
```

Server listens on `http://localhost:8080`.

## Project Structure

```
cmd/server/            entry point, dependency wiring, graceful shutdown
internal/model/        entities (Config, Timer)
internal/repository/   data access (JSON file storage)
internal/service/      business logic
internal/handler/      HTTP layer (Go 1.22+ ServeMux routing)
```

## Configuration

Pomodoro durations are stored in `data/pomodoro.json` (created on first save; defaults returned if the file does not exist).

```json
{
  "work_minutes": 25,
  "short_break_minutes": 5,
  "long_break_minutes": 15,
  "long_break_interval": 4
}
```

All values must be greater than `0`.

### Endpoints

```http
GET  /api/config
POST /api/config
```

Example:

```bash
curl http://localhost:8080/api/config

curl -X POST http://localhost:8080/api/config \
  -H 'Content-Type: application/json' \
  -d '{"work_minutes":50,"short_break_minutes":10,"long_break_minutes":20,"long_break_interval":3}'
```

## Timer

Timer state is held in memory and resets when the server restarts.

### Endpoints

```http
GET  /api/pomodoro
POST /api/pomodoro/start
POST /api/pomodoro/pause
POST /api/pomodoro/resume
POST /api/pomodoro/stop
```

`start` accepts an optional body:

```json
{"mode": "work"}
```

Modes: `work` (default), `short_break`, `long_break`.

Example:

```bash
curl -X POST http://localhost:8080/api/pomodoro/start \
  -H 'Content-Type: application/json' \
  -d '{"mode":"work"}'

curl http://localhost:8080/api/pomodoro

curl -X POST http://localhost:8080/api/pomodoro/pause
curl -X POST http://localhost:8080/api/pomodoro/resume
curl -X POST http://localhost:8080/api/pomodoro/stop
```

Status response:

```json
{
  "state": "running",
  "mode": "work",
  "remaining_seconds": 2999,
  "started_at": "2026-10-03T11:35:01+07:00"
}
```

States: `stopped`, `running`, `paused`.

Errors:

- `400` — invalid JSON, invalid config value, resume/stop in wrong state
- `409` — start while running, pause while not running

## Test

```bash
go vet ./...
go test ./...
```

## Notes

- Config storage: JSON file (swap `ConfigRepo` for a database later; it is already isolated behind the service layer).
- Timer state: in-memory only (planned upgrade: persist or broadcast over WebSocket).
