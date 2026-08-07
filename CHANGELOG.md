# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Added
- `domain.StorageService.Get(ctx, key) (io.ReadCloser, error)`, implemented by both `LocalStorage` and `S3Storage`, so downloads read from whichever backend actually stored the file.
- `domain.ErrObjectNotFound` sentinel, returned by `Get` on both backends (translated from `os.IsNotExist` locally and `*types.NoSuchKey` on S3), so handlers can distinguish "not found" from other storage failures without knowing the backend.
- `internal/config` package: a single `Config` struct loaded once in `main()` via `config.Load()`, replacing scattered `os.Getenv` calls and a hardcoded connection string. `AWS_S3_BUCKET_NAME` is required; `DATABASE_URL`, `AWS_REGION`, and `PORT` have defaults.

### Fixed
- **Downloads were broken for anything uploaded via `S3Storage`**: `HandleDownloadFile` read straight from the local `./uploads` directory regardless of storage backend. It now calls `storageService.Get` instead.
- **Path traversal via client-supplied filename**: `header.Filename` was embedded directly into the storage key; a filename containing `..` segments could write outside `LocalStorage`'s upload directory. Filenames are now passed through `filepath.Base` before use, and `LocalStorage` additionally rejects any resolved path outside its upload directory as defense in depth.
- `pgRepo.DB.Close()` was called twice (once via `defer`, once explicitly during graceful shutdown) — removed the redundant `defer`.
- `context.TODO()` used for loading AWS SDK config in `main()` replaced with `context.Background()`.
- `crypto/rand.Read`'s return values (including its error) were discarded entirely; the error is now checked and treated as a fatal request error.

### Changed
- `NewLocalStorage`/`NewS3Storage` now return `domain.StorageService` instead of their concrete types.

## 2026-08-06

### Added
- `S3Storage` implementation of `domain.StorageService`, using the AWS SDK v2 (`PutObject`/`DeleteObject`) and a `PresignClient` for generating presigned upload URLs (not yet wired into any handler).
- `.env.example` documenting `DATABASE_URL`, `AWS_REGION`, `AWS_S3_BUCKET_NAME`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.

### Changed
- `cmd/api/main.go` now wires `S3Storage` (not `LocalStorage`) for uploads, loading AWS config and the bucket name from the environment; the app exits at startup if `AWS_S3_BUCKET_NAME` is unset.
- `domain.StorageService` gained `GeneratePresignedUploadURL(ctx, key, expiration)`.

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
