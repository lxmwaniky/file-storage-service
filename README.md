# File Storage Service

A file upload, metadata tracking, and streaming service built in Go. Designed with clean-architecture layering, structured logging, and basic input validation.

---

## Architecture Notes

- **Core HTTP:** Go's native `net/http` and `ServeMux`, using native path wildcard routing (`GET /api/v1/files/{id}/download`), with configured read/write/idle timeouts and graceful shutdown on `SIGINT`/`SIGTERM`.
- **Storage split:** Raw file bytes are written and read through a `domain.StorageService` interface (`Save`/`Get`/`Delete`/`GeneratePresignedUploadURL`). Two implementations exist — `LocalStorage` (writes to `./uploads`) and `S3Storage` (S3 bucket, with presigned-URL generation). `cmd/api/main.go` currently wires up `S3Storage` only. Metadata (original name, size, MIME type, timestamp) is persisted to PostgreSQL via a single `INSERT`.
- **Observability:** Structured logging via the standard library `log/slog`, JSON-formatted.
- **Validation:** Upload size is capped with `http.MaxBytesReader` (10MB) to reject oversized payloads. Stored filenames are prefixed with a random hex identifier (16 bytes from `crypto/rand`) in the form `<hexid>_<basename>`, where the client-supplied filename is passed through `filepath.Base` (and `LocalStorage` additionally rejects any resolved path outside its upload directory) to prevent path traversal.
- **Config:** A single `internal/config.Config`, loaded once at startup in `main()`. `AWS_S3_BUCKET_NAME` is required (the app exits if it's unset); everything else falls back to a sane default.

Known gaps:
- `internal/usecase` exists as an empty placeholder package for future business-logic separation.
- The `DATABASE_URL` default (used when the env var isn't set) points at local dev credentials committed in `internal/config/config.go` — fine for local development, but should be required (no fallback) once there's a real deployment target.
- `GeneratePresignedUploadURL` is implemented on both storage backends but not yet wired into any HTTP handler.

---

## Project Structure

```text
.
├── cmd
│   └── api
│       └── main.go              # Entrypoint, config load, server bootstrap, graceful shutdown
├── internal
│   ├── config                   # Config struct + env loading
│   ├── delivery
│   │   └── http                 # HTTP router and handlers
│   ├── domain                   # Core models (File) and interfaces (StorageService)
│   ├── infrastructure
│   │   ├── repository           # PostgreSQL connection and data access
│   │   └── storage              # LocalStorage and S3Storage implementations of StorageService
│   └── usecase                  # Empty — planned business logic layer
├── uploads                      # Local storage for uploaded files (gitignored)
├── .env.example                 # Documented env vars
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

- Go (see `go.mod` for the exact version)
- A running PostgreSQL instance
- An S3 bucket and AWS credentials (the app currently requires S3 — see below)

### Environment variables

Copy `.env.example` and fill in values, then export them into your shell (there is no `.env` loader in the app yet):

| Variable | Required | Notes |
|---|---|---|
| `AWS_S3_BUCKET_NAME` | Yes | App exits at startup if unset. |
| `AWS_REGION` | No | Defaults to `af-south-1`. |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | Yes (or another AWS credential source) | Picked up by the default AWS SDK credential chain. |
| `DATABASE_URL` | No | Falls back to a local dev Postgres URL — see Known gaps. |
| `PORT` | No | Defaults to `:8080`. |

### Database setup

Create the `files` table:

```sql
CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(64) PRIMARY KEY,
    original_name TEXT NOT NULL,
    stored_filename VARCHAR(128) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(128) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Running the service

```bash
go run ./cmd/api
```

The server listens on `:8080` (or `$PORT`) and shuts down gracefully on Ctrl+C.

---

## API

### `GET /health`

Health check.

```bash
curl http://localhost:8080/health
```

### `POST /api/v1/upload`

Accepts a multipart form file upload (field name `file`), rejects payloads over 10MB, saves the content via the storage service (S3, as currently wired), and writes the metadata row to PostgreSQL.

```bash
curl -F "file=@path/to/file.pdf" http://localhost:8080/api/v1/upload
```

### `GET /api/v1/files`

Returns metadata for all uploaded files.

```bash
curl http://localhost:8080/api/v1/files
```

### `GET /api/v1/files/{id}/download`

Looks up file metadata by ID, fetches the object from the storage service, and streams it back with its original filename and MIME type.

```bash
curl -O -J http://localhost:8080/api/v1/files/{id}/download
```
