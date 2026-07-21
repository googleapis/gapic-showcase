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
- **`query`**: Queries current session status and returns `X-Goog-Upload-Status` and `X-Goog-Upload-Size-Received` (`200 OK` for active/final sessions; `410 Gone` for cancelled sessions).
- **`cancel`**: Cancels the session (`X-Goog-Upload-Status: cancelled`).

---

## 3. Test Scenarios & Failure Injection

The middleware supports injecting failure scenarios via HTTP headers for testing client retry and error-recovery behavior:

- **`X-Goog-Test-Scenario`**: Specifies the scenario (`non_fatal_error_on_start`, `fatal_error_on_start`, `non_fatal_error_on_chunk_upload`, `non_fatal_error_on_query`, `chunk_granularity`).
- **`X-Goog-Test-Scenario-Config`**: JSON configuration string controlling the injected error (`client_uuid`, `error_code`, `failure_count`, `after_offset`, `action_after_failures`).

### Start-Call Scenario Isolation (`client_uuid`)
Because `start` commands occur before a session ID (`sid`) is generated, failure counts for `non_fatal_error_on_start` are tracked per `client_uuid` (or remote address if omitted). Passing a unique `client_uuid` in `X-Goog-Test-Scenario-Config` ensures concurrent test executions do not collide:

```json
{
  "client_uuid": "test-client-run-abc",
  "error_code": 503,
  "failure_count": 1
}
```
