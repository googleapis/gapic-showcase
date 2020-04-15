# Release History

### v0.8.1 / 2020-04-15
- Fix bug in mTLS configuration resolution

### v0.8.0 / 2020-04-14
- Add mtls support with user provided cert/key to server
- Regen cli:
  - Paginated RPCs only collect a single page
  - Default page_size changed from 0 to 10 to avoid short circuiting
- Regen client:
  - clientHook support added
- Update CI use of go.mod
- Update dependencies

### v0.7.0 / 2020-02-06
- Regen client and protobuf code for Go grpc.ClientConn interface
- Update dependencies

### v0.6.1 / 2019-11-01
- Fix the resource name for Blurb.

### v0.6.0 / 2019-11-01
- Add a gRPC ServiceConfig for microgenerator retry config
- Regen client code with retry config
- Update dependencies

### v0.5.0 / 2019-09-04
- Update to Go version 1.13
- Update dependencies
- Add trailers testing support to Echo
- Fix pagination in operations service

### v0.4.0 / 2019-08-13
- Add dummy LRO service
- Dependency updates

### v0.3.0 / 2019-08-09
- Remove nodejs server implementation
- Update dependencies
- Update golang docker tag to v1.12
- Add Block method to Echo service
- Enable kotlin smoke test
- Add renovate.json

### v0.2.4 / 2019-07-11
- Update `grpc-fallback-go` version to `v0.1.3`

### v0.2.3 / 2019-07-09
- Update `grpc-fallback-go` version to `v0.1.2`

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
