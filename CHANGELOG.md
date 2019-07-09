# Release History

### v0.2.2 / 2019-07-09

- Update `grpc-fallback-go` version to `v0.1.1`

### v0.2.1 / 2019-07-03
- Add fallback-proxy to `gapic-showcase run` via grpc-fallback-go
- Expose fallback-proxy port in Dockerfile
- Tidy `go.mod`

### v0.2.0 / 2019-05-24
- Regenerate GAPIC & GCLI with small updates
- Update resource annotations
- Fix bug in README
- Add Node.js EchoService implementation
- Remove extraneous logging

### v0.1.1 / 2019-04-17
- Regenerate GAPIC & GCLI to capture fixes for paged RPCs & commands

### v0.1.0 / 2019-04-04
- Beta release.

### v0.0.16 / 2019-03-25
- Fixing some field names in path templates

### v0.0.15 / 2019-03-25
- Fixing path templates to make sure curly braces match
- Use Go modules

### v0.0.14 / 2019-03-25
- Serve Testing service CLI service
- Ensure all path templates start with `/`

### v0.0.13 / 2019-03-01
- Fix issue which tombstones users.

### v0.0.12 / 2019-02-20
- Remove google.api.client_package proto annotations.

### v0.0.11 / 2019-02-19
- Update GAPIC config proto annotations.

### v0.0.10 / 2019-01-29
- Expose messaging and identity services when running `gapic-showcase run`.
- Refactor `Echo.WaitRequest` to follow API style for denoting time to live.
- Use GCLI Generated Code for the CLI cmd.
