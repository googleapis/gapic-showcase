# GAPIC Showcase

[![Release Level][releaselevelimg]][releaselevel]
[![CircleCI][circleimg]][circle]
[![Code Coverage][codecovimg]][codecov]
[![GoDoc][godocimg]][godoc]

> GAPIC (Generated API Client) Showcase is an API that showcases features used
by GAPICs to make calling Google APIs an enjoyable experience.

The main goal of GAPIC Showcase is to be a testing tool that will be able to
verify the features that a GAPIC implements. GAPIC Showcase services aim to be
representative of all API Client configurations of Google APIs. With this aim in
mind, gapic-generators that pass gapic-showcase can have reasonable confidence
in the clients they are generating.

## Services
The services of GAPIC Showcase API can be found in [schema/](schema/). Please
note that these protocol buffer files are not able to be compiled in isolation.
To get the services staged alongside their dependencies, please see check out
our [releases](https://github.com/googleapis/gapic-showcase/releases) page.

## GAPIC Showcase CLI Tool
### Installation
The GAPIC Showcase CLI Tool can be installed using three different mechanisms,
downloading the compiled binary from our our [releases](https://github.com/googleapis/gapic-showcase/releases)
page, pulling our released docker image from [google container registry](https://gcr.io/gapic-showcase/gapic-showcase),
or simply by using go commands.

#### Binary Installation
```sh
$ export GAPIC_SHOWCASE_VERSION=0.0.10
$ export OS=linux
$ export ARCH=amd64
$ curl -L https://github.com/googleapis/gapic-showcase/releases/download/v${GAPIC_SHOWCASE_VERSION}/gapic-showcase-${GAPIC_SHOWCASE_VERSION}-${OS}-${ARCH} | sudo tar -zx -- --directory /usr/local/bin/
$ gapic-showcase run
> 2018/09/19 02:13:09 Showcase listening on port: :7469
```

#### Docker Installation
```sh
$ export GAPIC_SHOWCASE_VERSION=0.0.10
$ docker pull gcr.io/gapic-showcase/gapic-showcase:${GAPIC_SHOWCASE_VERSION}
$ docker run -it gcr.io/gapic-showcase/gapic-showcase:${GAPIC_SHOWCASE_VERSION}
> 2018/09/19 02:13:09 Showcase listening on port: :7469
```

#### Go Installation
```sh
$ go install github.com/googleapis/gapic-showcase
$ gapic-showcase run
> 2018/09/19 02:13:09 Showcase listening on port: :7469
```

<!---
TODO(landrito): figure out a blessed way to install by version using go
commands.
-->
_* Bear in mind this is not a versioned installation so no versioning guarantees
hold using this installation method._

## Example Usage - Implementing a GAPIC Showcase Integration Test

### Step 1. Generate a gapic-showcase client
To start, a user will download the gapic-showcase protobuf files or proto
descriptor set from a gapic-showcase release. The user will then feed these
protobuf files into their gapic-generator. This client will be the client used
for integration testing their gapic- generator.

```sh
$ export GAPIC_SHOWCASE_VERSION=0.0.10
$ curl -L https://github.com/googleapis/gapic-showcase/releases/download/v${GAPIC_SHOWCASE_VERSION}/gapic-showcase-${GAPIC_SHOWCASE_VERSION}-protos.tar.gz | sudo tar -zx
$ protoc google/showcase/v1alpha3/*.proto \
    --proto_path=. \
    --${YOUR_GAPIC_GENERATOR}_out=/dest/
```

### Step 2. Write Integration Tests
Write an Integration test which calls the gapic-showcase server.

#### Ruby Example for [gapic-generator](https://github.com/googleapis/gapic-generator)
<!---
TODO(landrito): Add testing service stuff once it is implemented.
-->
```rb
describe Google::Showcase::V1alpha1::EchoClient do
  before(:all) do
    # gapic-showcase does not implement any auth so an insecure channel must be
    # used.
    channel = credentials: GRPC::Core::Channel.new(
      "localhost:7469", nil, :this_channel_is_insecure)
    @echo_client = Google::Showcase::V1alpha2::EchoClient.new(channel)
  end

  describe 'echo' do
    it 'invokes echo without error' do
      # Create expected grpc response
      content = "Echo Content"
      expected_response = { content: content }
      expected_response = Google::Gax::to_proto(
      expected_response, Google::Showcase::V1alpha1::EchoResponse)

      # Call method
      response = @echo_client.echo(content: content)

      # Verify the response
      assert_equal(expected_response, response)
    end
  end
end
```

### Step 3. Run the Showcase Server
The integration test needs a server to send its requests to. Download and run
the server so that gapic-showcase is available for the tests.

```sh
$ export GAPIC_SHOWCASE_VERSION=0.0.10
$ export OS=linux
$ export ARCH=amd64
$ curl -L https://github.com/googleapis/gapic-showcase/releases/download/v${GAPIC_SHOWCASE_VERSION}/gapic-showcase-${OS}-${ARCH} | sudo tar -zx -- --directory /usr/local/bin/
$ gapic-showcase run
> 2018/09/19 02:13:09 Showcase listening on port: :7469
```

### Step 4. Run integration tests against the server
Now the integration test is ready to be run. Invoke your test!

#### Ruby Example for [gapic-generator](https://github.com/googleapis/gapic-generator)
```sh
$ bundle install
$ bundle exec ruby gapic-showcase-integration-test.rb
```
<!---
TODO(landrito): Add test report once it is implemented.
-->

## Released Artifacts
GAPIC Showcase releases three main artifacts, a CLI tool, the gapic-showcase
service protobuf files staged alongside its dependencies, and a protocol buffer
descriptor set compiled from the gapic-showcase service protos.

Check out our [releases](https://github.com/googleapis/gapic-showcase/releases) page to see our released artifacts.

### CLI Tool
The GAPIC Showcase CLI Tool is used for two purposes, to run the
gapic-showcase server, and to make requests to an already running gapic-showcase
server.

Generally, any questions about using the CLI tool can be answered by running the
CLI tool with the `--help` flag which will supply usage documentation.

```sh
$ gapic-showcase [command?] --help
```

#### Starting the Server
The primary purpose of the CLI tool will be starting the showcase server. This
server will expose the GAPIC Showcase services port 7469 by default.

##### Spinning Up the Server
```sh
$ gapic-showcase run
> 2018/09/19 01:57:09 Showcase listening on port: :7469

$ gapic-showcase run --port 1234
> 2018/09/19 01:57:09 Showcase listening on port: :1234
```

#### Making Showcase Service Calls
The CLI tool will also be able to make requests to a running showcase server.
This allows you to have a simple way to interact and tinker with the Showcase
API.

##### Example
```sh
$ gapic-showcase run                    
> 2018/09/19 02:13:09 Showcase listening on port: :7469
> 2018/09/19 02:14:08 Received Unary Request for Method: /google.showcase.v1alpha3.Echo/Echo
> 2018/09/19 02:14:08     Request:  content:"hello world"
> 2018/09/19 02:14:08     Returning Response: content:"hello world"
```
```sh
$ gapic-showcase echo hello world
> 2018/09/19 02:14:08 Sent Request: content: "hello world"
> 2018/09/19 02:14:08 Got Response: content: "hello world"
```

### Staged Protocol Buffer Files
The [proto files](schema/) found in the gapic-showcase repository are not compilable in
isolation. This is to avoid duplication of the protofiles that showcase depends
on, namely the API and API client configurations found in the `input-contract` branch of
[api-common-protos](https://github.com/googleapis/api-common-protos/tree/input-contract).
To give the user everything that is needed to compile the showcase protocol
buffer files, every [release](https://github.com/googleapis/gapic-showcase/releases)
will have attached a tarball containing a snapshot of the gapic-showcase
protocol buffer files staged alongside their dependencies.

### Compiled Proto Descriptors
The compiled proto descriptors for the staged gapic-showcase protos discussed in
the previous section will be attached to every release. This is intended to make
it easier to generate clients removing the necessary step of installing protoc.

## Versioning
GAPIC Showcase follows semantic versioning in which all artifacts that are
released for a certain version are guaranteed to be compatible with one another.
To be more explicit, for a certain version, the interfaces and types exposed by
the protobuf files are compatible with the interface of the implemented server.
Users of gapic-showcase are expected to implement integration tests against a
certain version of gapic-showcase rather than implementing against
gapic-showcase at head.

## Supported Go Versions
GAPIC Showcase is supported for go versions 1.11 and later.

## FAQ

### Is this Showcase API publicly served?

This API is not publicly served. Users of gapic-showcase are expected to run the
server locally.

## Disclaimer

This is not an official Google product.

[circle]: https://circleci.com/gh/googleapis/gapic-showcase
[circleimg]: https://circleci.com/gh/googleapis/gapic-showcase.svg?style=shield
[codecovimg]: https://codecov.io/github/googleapis/gapic-showcase/coverage.svg?branch=master
[codecov]: https://codecov.io/github/googleapis/gapic-showcase?branch=master
[godoc]: https://godoc.org/github.com/googleapis/gapic-showcase/server
[godocimg]: https://godoc.org/github.com/googleapis/gapic-showcase/server?status.svg
[releaselevel]: https://cloud.google.com/terms/launch-stages
[releaselevelimg]: https://img.shields.io/badge/release%20level-alpha-red.svg?style&#x3D;flat
