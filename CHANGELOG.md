# Release History

## [0.27.0](https://github.com/googleapis/gapic-showcase/compare/v0.26.1...v0.27.0) (2023-04-19)


### Features

* Add iam and location mixin rest handlers ([#1300](https://github.com/googleapis/gapic-showcase/issues/1300)) ([6adab7b](https://github.com/googleapis/gapic-showcase/commit/6adab7bb4e7c979f90441a37cde36ecda0bee68f))
* Add proto + logic for streaming sequence ([#1266](https://github.com/googleapis/gapic-showcase/issues/1266)) ([82814d8](https://github.com/googleapis/gapic-showcase/commit/82814d8a8ada26bf9a831497b667183edf1f0e4f))


### Bug Fixes

* Handle already quoted protojson wkt ([#1294](https://github.com/googleapis/gapic-showcase/issues/1294)) ([f74f03d](https://github.com/googleapis/gapic-showcase/commit/f74f03d50ea8aa5832f830a6f92e16e99fb3e623))

## [0.26.1](https://github.com/googleapis/gapic-showcase/compare/v0.26.0...v0.26.1) (2023-03-28)


### Bug Fixes

* **rest:** Properly handle string-encoded well-known types in URLs ([#1282](https://github.com/googleapis/gapic-showcase/issues/1282)) ([579fe72](https://github.com/googleapis/gapic-showcase/commit/579fe729c0b506cdf48cd834beb14a8ffb4b5994))

## [0.26.0](https://github.com/googleapis/gapic-showcase/compare/v0.25.0...v0.26.0) (2023-03-07)


### Features

* **go:** Update Go version to 1.19 ([#1225](https://github.com/googleapis/gapic-showcase/issues/1225)) ([d4b108e](https://github.com/googleapis/gapic-showcase/commit/d4b108e16dc91c0ea6d4dec3dca4d3270d3bf47a))


### Bug Fixes

* Build assets for darwin/arm64 ([#1267](https://github.com/googleapis/gapic-showcase/issues/1267)) ([0833a57](https://github.com/googleapis/gapic-showcase/commit/0833a579131c14582b053f26698fdfe93e465d87))
* Export showcase_v1beta1.yaml from BUILD.bazel to support external GAPIC generation ([#1223](https://github.com/googleapis/gapic-showcase/issues/1223)) ([5076348](https://github.com/googleapis/gapic-showcase/commit/507634898e208b7ff88784e4ec5f0efd22bff9ab))
* Handle x-http-method-override for PATCH as POST ([#1262](https://github.com/googleapis/gapic-showcase/issues/1262)) ([4070ce3](https://github.com/googleapis/gapic-showcase/commit/4070ce331bd5e852ccb2f4f2267dce80a9dda9c4))
* Use quotes around extreme int64 values ([#1206](https://github.com/googleapis/gapic-showcase/issues/1206)) ([c9d9ff1](https://github.com/googleapis/gapic-showcase/commit/c9d9ff191bfd72fe8563625be4074fe4659585d6)), closes [#1205](https://github.com/googleapis/gapic-showcase/issues/1205)

## [0.25.0](https://github.com/googleapis/gapic-showcase/compare/v0.24.0...v0.25.0) (2022-09-01)


### Features

* Support FieldMask in Updates ([#1197](https://github.com/googleapis/gapic-showcase/issues/1197)) ([cdb4ce6](https://github.com/googleapis/gapic-showcase/commit/cdb4ce63778a8ea6bcef0006c1e4c4a50da45a6c))


### Bug Fixes

* Use resource field in http body for Updates ([#1198](https://github.com/googleapis/gapic-showcase/issues/1198)) ([48a2632](https://github.com/googleapis/gapic-showcase/commit/48a2632b5f25af24b239216f1ff079ed60de2a61))

## [0.24.0](https://github.com/googleapis/gapic-showcase/compare/v0.23.0...v0.24.0) (2022-07-28)


### Features

* **regapic:** accept numeric enums, allow testing enum round trips ([#1159](https://github.com/googleapis/gapic-showcase/issues/1159)) ([1e863ae](https://github.com/googleapis/gapic-showcase/commit/1e863ae834ad58453c1d74cf593be1188ff15033))

## [0.23.0](https://github.com/googleapis/gapic-showcase/compare/v0.22.0...v0.23.0) (2022-07-27)


### Features

* add binding testing and multiple bindings test data to compliance service ([#1150](https://github.com/googleapis/gapic-showcase/issues/1150)) ([9d43ed0](https://github.com/googleapis/gapic-showcase/commit/9d43ed0621e4c9549a3d92a53dedc5dd57e9bec2))


### Bug Fixes

* remove broken test cases in compliance ([#1151](https://github.com/googleapis/gapic-showcase/issues/1151)) ([a56df9a](https://github.com/googleapis/gapic-showcase/commit/a56df9ab5ad9cdd6cb0385ba9c74263b95538ad8))

## [0.22.0](https://github.com/googleapis/gapic-showcase/compare/v0.21.0...v0.22.0) (2022-06-13)


### Features

* Support LRO mixins over REST ([#1118](https://github.com/googleapis/gapic-showcase/issues/1118)) ([5ca6fe1](https://github.com/googleapis/gapic-showcase/commit/5ca6fe1b8ea7e5645e87718e84b1198ad8ce9c63))

## [0.21.0](https://github.com/googleapis/gapic-showcase/compare/v0.20.0...v0.21.0) (2022-06-08)


### Features

* respond to requests specifying response enum values be JSON-encoded as ints ([#1111](https://github.com/googleapis/gapic-showcase/issues/1111)) ([5389bd1](https://github.com/googleapis/gapic-showcase/commit/5389bd17aedb7f0c8a8de562421222e703589823))


### Bug Fixes

* **genrest:** pass http request context to service handler ([#1088](https://github.com/googleapis/gapic-showcase/issues/1088)) ([bad9b6b](https://github.com/googleapis/gapic-showcase/commit/bad9b6b89b0f75d3ec8408610d329068775703e5))

## [0.20.0](https://github.com/googleapis/gapic-showcase/compare/v0.19.5...v0.20.0) (2022-05-10)


### Features

* **genrest:** format rest errors as Google Errors ([#1082](https://github.com/googleapis/gapic-showcase/issues/1082)) ([e49f134](https://github.com/googleapis/gapic-showcase/commit/e49f134dc54ae734c716f8649fc3279efd68916f))

### [0.19.5](https://github.com/googleapis/gapic-showcase/compare/v0.19.4...v0.19.5) (2022-03-08)


### Bug Fixes

* add a routing.proto dependency to showcase proto bazel target ([#1033](https://github.com/googleapis/gapic-showcase/issues/1033)) ([e2aa303](https://github.com/googleapis/gapic-showcase/commit/e2aa3033e67c0c50c5650512211001e1ad29d36a))

### [0.19.4](https://github.com/googleapis/gapic-showcase/compare/v0.19.3...v0.19.4) (2022-03-03)


### Bug Fixes

* **ci:** yet another attempt to fix asset ci ([#1026](https://github.com/googleapis/gapic-showcase/issues/1026)) ([32a2603](https://github.com/googleapis/gapic-showcase/commit/32a2603550220dcd0a86d533a883fa1ba860a3b1))

### [0.19.3](https://github.com/googleapis/gapic-showcase/compare/v0.19.2...v0.19.3) (2022-03-03)


### Bug Fixes

* **ci:** asset version shouldn't include the leading v ([#1023](https://github.com/googleapis/gapic-showcase/issues/1023)) ([a78d624](https://github.com/googleapis/gapic-showcase/commit/a78d6243b9cbfa088b6863954b0bf6dec91c6786))

### [0.19.2](https://github.com/googleapis/gapic-showcase/compare/v0.19.1...v0.19.2) (2022-03-02)


### Bug Fixes

* use tag_name for proto-assets upload ([#1021](https://github.com/googleapis/gapic-showcase/issues/1021)) ([c3f2f36](https://github.com/googleapis/gapic-showcase/commit/c3f2f36128a3c7b5c51fc1bd87c83f8f833adf9a))

### [0.19.1](https://github.com/googleapis/gapic-showcase/compare/v0.19.0...v0.19.1) (2022-03-02)


### Bug Fixes

* ruby package name ([#1015](https://github.com/googleapis/gapic-showcase/issues/1015)) ([7289f2e](https://github.com/googleapis/gapic-showcase/commit/7289f2e20afa198695d3da9d5a5362a161032244))

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
