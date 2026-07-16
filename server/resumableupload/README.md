# `ResumableUploadService` in `gapic-showcase`: Methods & Test Scenarios

`gapic-showcase` provides a canonical protobuf service (`google.showcase.v1beta1.ResumableUploadService`) and HTTP middleware (`server/resumableupload`) implementing the **Scotty Universal Resumable Upload Protocol**.

---

## 1. RPC Methods

### `UploadMedia`
- **Proto RPC**: `rpc UploadMedia(UploadMediaRequest) returns (UploadMediaResponse)`
- **HTTP Path**: `POST /resumable/upload/v1beta1/files:upload`
- **Purpose**: Canonical RPC method annotated with `google.api.http` enabling resumable upload session creation.

---

## 2. Test Scenarios (`X-Goog-Test-Scenario`)

All test scenarios are controlled by passing custom HTTP headers (`X-Goog-Test-Scenario` and optional `X-Goog-Test-Scenario-Config`) when initiating the session (`X-Goog-Upload-Command: start`).

| Scenario Name | Description | Header Configuration | Expected Client Behavior |
| :--- | :--- | :--- | :--- |
| **Happy Path** | Default upload flow. Returns `X-Goog-Upload-Chunk-Granularity: 1` and completes upload. | None | Completes upload successfully on first attempt. |
| **`non_fatal_error_on_start`** | Returns a retriable HTTP error (`503 Service Unavailable` by default) on initial `start` command. | `X-Goog-Test-Scenario: non_fatal_error_on_start`<br/>`X-Goog-Test-Scenario-Config: {"error_code": 503, "failure_count": 1}` | Automatically retries initial `start` request and succeeds. |
| **`fatal_error_on_start`** | Returns a fatal HTTP error (`403 Forbidden` by default) on initial `start` command. | `X-Goog-Test-Scenario: fatal_error_on_start`<br/>`X-Goog-Test-Scenario-Config: {"error_code": 403}` | Aborts session and raises an error with status code `403`. |
| **`non_fatal_error_on_chunk_upload`** | Returns a retriable error (`503`) during chunk upload, triggering recovery query exchange. | `X-Goog-Test-Scenario: non_fatal_error_on_chunk_upload`<br/>`X-Goog-Test-Scenario-Config: {"error_code": 503, "failure_count": 1, "after_offset": 1024}` | Enters `RECOVERY` phase, queries committed offset (`X-Goog-Upload-Command: query`), and resumes uploading remaining bytes. |
| **`non_fatal_error_on_query`** | Returns a retriable error (`503`) on the recovery `query` command. | `X-Goog-Test-Scenario: non_fatal_error_on_query`<br/>`X-Goog-Test-Scenario-Config: {"error_code": 503, "failure_count": 1}` | Retries the `query` request until server returns `X-Goog-Upload-Size-Received`, then resumes transmission. |
| **`chunk_granularity`** | Returns `X-Goog-Upload-Chunk-Granularity: 256` on `start` and rejects non-final chunks whose length is not a multiple of `256`. | `X-Goog-Test-Scenario: chunk_granularity` | Parses granularity from `start` response and rounds chunk boundaries to multiples of `256` bytes. |
