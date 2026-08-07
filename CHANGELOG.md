# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Added
- `S3Storage` implementation of `domain.StorageService`, using the AWS SDK v2 (`PutObject`/`DeleteObject`) and a `PresignClient` for generating presigned upload URLs (not yet wired into any handler).
- `.env.example` documenting `DATABASE_URL`, `AWS_REGION`, `AWS_S3_BUCKET_NAME`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.

### Changed
- `cmd/api/main.go` now wires `S3Storage` (not `LocalStorage`) for uploads, loading AWS config and the bucket name from the environment; the app exits at startup if `AWS_S3_BUCKET_NAME` is unset.
- `domain.StorageService` gained `GeneratePresignedUploadURL(ctx, key, expiration)`.

### Known issues
- `GET /api/v1/files/{id}/download` still reads directly from the local `./uploads` directory, so files uploaded via `S3Storage` cannot be downloaded through the service yet.
- `DATABASE_URL` is documented in `.env.example` but not read anywhere — the Postgres connection string is still hardcoded in `main.go`.

## 2026-08-05

### Added
- Pluggable `domain.StorageService` interface with a `LocalStorage` implementation, decoupling upload/download handlers from the local filesystem.
- Graceful shutdown on `SIGINT`/`SIGTERM`, with a configured `http.Server` (read/write/idle timeouts) and a shutdown timeout for in-flight requests.

### Changed
- Stored filenames now follow `<random-hex-id>_<original-filename>` instead of `<random-hex-id><ext>`, so the original filename is retained on disk.
- Upload and download handlers log more context around failures (e.g. missing file ID, DB record with no matching file on disk).

## 2026-08-05 (earlier)

### Added
- `GET /api/v1/files` — list metadata for all uploaded files.
- `GET /api/v1/files/{id}/download` — stream a file back by ID with its original filename and MIME type.
- PostgreSQL-backed file metadata storage (`FileRepository`, `files` table).

### Changed
- Upload handler persists file metadata (original name, stored name, size, MIME type, timestamp) to PostgreSQL after saving the file.

## 2026-08-04

### Added
- Initial HTTP server with `GET /health` and `POST /api/v1/upload`.
- Upload validation: 10MB size cap via `http.MaxBytesReader`, randomized stored filenames.
- `File` domain model and initial PostgreSQL repository implementation.
- Project `.gitignore` and initial README.
