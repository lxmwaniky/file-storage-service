# File Storage Service

A minimal HTTP file storage service written in Go.

## Run

```
go run ./cmd/api
```

Server starts on `:8080`.

## Endpoints

### `GET /health`

Health check.

```
curl http://localhost:8080/health
```

### `POST /api/v1/upload`

Uploads a file (multipart form, field name `file`). Max size 10MB. Files are stored under `./uploads` with a randomized filename.

```
curl -F "file=@example.jpg" http://localhost:8080/api/v1/upload
```
