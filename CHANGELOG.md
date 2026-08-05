# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Added
- Pluggable `domain.StorageService` interface with a `LocalStorage` implementation, decoupling upload/download handlers from the local filesystem.
- Graceful shutdown on `SIGINT`/`SIGTERM`, with a configured `http.Server` (read/write/idle timeouts) and a shutdown timeout for in-flight requests.

### Changed
- Stored filenames now follow `<random-hex-id>_<original-filename>` instead of `<random-hex-id><ext>`, so the original filename is retained on disk.
- Upload and download handlers log more context around failures (e.g. missing file ID, DB record with no matching file on disk).

## 2026-08-05

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
