# Release History

### v0.19.0 / 2022-01-31
- update api-common-protos submodule
- add ability to echo headers and added several routing annotations to the `Echo` method
- enable generation of both grpc and rest clients

### v0.18.0 / 2021-12-06
- add `parent` to method signature for `Messaging.SearchBlurbs()`
- update `RELEASING.md` instructions

### v0.17.0 / 2021-11-02
- Implement server streaming RPCs over REST, using chunked encoding.
- Implement RPCs that map to PUT and PATCH HTTP verbs
- Check that REST RPCs using HTTP GET or DELETE don't contain bodies.
- Disable TypeScript smoke tests pending upstream fixes (TS generator Docker image).

### v0.16.0 / 2021-06-16
- Require incoming REST requests to have expected `x-goog-api-client` header tokens
- Allow mTLS to work over gRPC when using `cmux` to also listen to REST requests on the same port
- Make REST `PATCH` methods work
- Fix multi-line truncation in release notes
- Add Docker push instructions to RELEASING.md

### v0.15.0 / 2021-05-05
- Enforce `Content-Type: application/json` in the bodies of REST requests
- Enforce correct `optional` field presence/absence in test suite requests (bodies and query strings)
- Lower-camel-case field names in `compliance_suite.json`

### v0.14.0 / 2021-04-27
- Fix collision between operation helper for `Echo.Wait` and generated mixin `Operations.WaitOpertation`
- REST endpoints: ensure enum values are received as string values
- REST endpoints: ensure full body responses
- Rest endpoints: enforce lower-camel-cased field names in request bodies and query params
- fix windows binary upload
- fix go vet/lint warnings
- pin Go version in CI
- fix release asset version
- add Code of Conduct
- add `SECURITY.md`

### v0.13.0 / 2021-02-24
- Auto-generate REST endpoints for Showcase services via `genrest` (partial)
- Add Compliance service for generators to use to test REST-transcoding their
  protos and RPCs (partial)
- Add mix-in service implementations
- Update API Service config with mix-ins and more
- Add Bazel proto_library targets for schema/
- Migrated to GitHub Actions
- Regen client & CLI with small updates
- Update dependencies

### v0.12.0 / 2020-08-12
- Add client-side retry/deadline testing surface
- Regen client & CLI with small updates
- Update dependencies

### v0.11.0 / 2020-05-26
- Add non-slash resource name patterns to Blurb resource
- Fix typo in User-parented Blurb resource patterns
- Add an enum to EchoRequest/EchoResponse
- Regen CLI with new fields
- Update dependencies

### v0.10.1 / 2020-05-20
- Fix UpdateUser handler response to send entire updated resource
- Note: non-slash resource name changes are not included in this release

### v0.10.0 / 2020-05-18
- Add use of proto3_optional in schema
- Upgrade CI protoc to v3.12.0
- Regen CLI with new proto3_optional fields
- Update dependencies

### v0.9.0 / 2020-04-22
- Print gRPC request headers in verbose mode (`gapic-showcase run -v`)
- Add TypeScript smoke tests
- Fix Kotlin smoke tests

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
