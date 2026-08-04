# File Storage Service

A file upload, metadata tracking, and streaming service built in Go. Designed with a clean-architecture layering, structured logging, memory-safe streaming, and basic input validation.

---

## Architecture Notes

- **Core HTTP:** Go's native `net/http` and `ServeMux`, using native path wildcard routing (`GET /api/v1/files/{id}/download`).
- **Storage split:** Raw file bytes are streamed to local disk (`./uploads`); metadata (original name, size, MIME type, timestamp) is persisted to PostgreSQL via a single `INSERT`.
- **Observability:** Structured logging via the standard library `log/slog`, JSON-formatted.
- **Validation:** Upload size is capped with `http.MaxBytesReader` (10MB) to reject oversized payloads. Stored filenames are replaced with randomized hex-encoded identifiers (16 bytes from `crypto/rand`) rather than the client-supplied name, avoiding path traversal via user input.
- **Streaming:** Uploads and downloads use `io.Copy`, so file contents aren't buffered fully in memory.

Known gaps: the PostgreSQL connection string is currently hardcoded in `cmd/api/main.go` rather than read from configuration; `internal/config`, `internal/usecase`, and `internal/infrastructure/storage` exist as empty placeholder packages for future work.

---

## Project Structure

```text
.
├── cmd
│   └── api
│       └── main.go              # Application entrypoint & server bootstrap
├── internal
│   ├── config                   # Empty — planned config loading
│   ├── delivery
│   │   └── http                 # HTTP router and handlers
│   ├── domain                   # Core models (File)
│   ├── infrastructure
│   │   ├── repository           # PostgreSQL connection and data access
│   │   └── storage              # Empty — planned storage abstraction
│   └── usecase                  # Empty — planned business logic layer
├── uploads                      # Local storage for uploaded files (gitignored)
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

- Go (see `go.mod` for the exact version)
- A running PostgreSQL instance

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

The connection string used by the app is set in [cmd/api/main.go](cmd/api/main.go) — update it to match your local database.

### Running the service

```bash
go run ./cmd/api
```

The server listens on `:8080`.

---

## API

### `GET /health`

Health check.

```bash
curl http://localhost:8080/health
```

### `POST /api/v1/upload`

Accepts a multipart form file upload (field name `file`), rejects payloads over 10MB, generates a random stored filename, streams the content to disk, and writes the metadata row to PostgreSQL.

```bash
curl -F "file=@path/to/file.pdf" http://localhost:8080/api/v1/upload
```

### `GET /api/v1/files`

Returns metadata for all uploaded files.

```bash
curl http://localhost:8080/api/v1/files
```

### `GET /api/v1/files/{id}/download`

Looks up file metadata by ID, checks the file exists on disk, and streams it back with its original filename and MIME type.

```bash
curl -O -J http://localhost:8080/api/v1/files/{id}/download
```
