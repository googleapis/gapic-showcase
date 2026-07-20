# `ResumableUploadService` in `gapic-showcase`

`gapic-showcase` provides a canonical protobuf service (`google.showcase.v1beta1.ResumableUploadService`) and HTTP middleware (`server/resumableupload`) implementing the core **Scotty Universal Resumable Upload Protocol**.

---

## 1. RPC Methods

### `UploadMedia`
- **Proto RPC**: `rpc UploadMedia(UploadMediaRequest) returns (UploadMediaResponse)`
- **HTTP Path**: `POST /resumable/upload/v1beta1/files:upload`
- **Purpose**: Canonical RPC method annotated with `google.api.http` enabling resumable upload session creation.

---

## 2. Core Resumable Upload Protocol Commands Handled

The middleware inspects `X-Goog-Upload-Command` and implements the core session exchange:

- **`start`**: Initiates a session. Returns `200 OK` with `X-Goog-Upload-Status: active`, `X-Goog-Upload-URL`, `X-Goog-Upload-Control-URL`, and `X-Goog-Upload-Chunk-Granularity: 262144`.
- **`upload`**: Appends chunk data to the session buffer after validating `X-Goog-Upload-Offset`.
- **`upload, finalize`**: Commits the final chunk, sets session status to `final`, and returns the JSON response `{"name":"<id>","size":<total_bytes>}`.
- **`query`**: Queries current session status and returns `X-Goog-Upload-Status` and `X-Goog-Upload-Size-Received`.
- **`cancel`**: Cancels the session (`X-Goog-Upload-Status: cancelled`).
