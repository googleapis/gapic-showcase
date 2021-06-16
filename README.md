# GAPIC Showcase

[![Release Level][releaselevelimg]][releaselevel]
[![Code Coverage][codecovimg]][codecov]
[![GoDoc][godocimg]][godoc]

> GAPIC (Generated API Client) Showcase: an API to test your gRPC clients.

GAPIC Showcase is an API to allow you to test a gRPC server and client.
It is a simple showcase of an API that follows the [Cloud APIs design guide](https://cloud.google.com/apis/design/).
This project provides a server and client implementation of the API that
can be run locally.

## Installation
The GAPIC Showcase CLI can be installed using three different mechanisms:
pulling a docker image from [Google Container Registry](https://gcr.io/gapic-images/gapic-showcase),
downloading the compiled binary from our our [releases](https://github.com/googleapis/gapic-showcase/releases)
page, or simply by installing from source using go.

### Docker
```sh
$ docker pull gcr.io/gapic-images/gapic-showcase:0.16.0
$ docker run \
    --rm \
    -p 7469:7469/tcp \
    -p 7469:7469/udp \
    gcr.io/gapic-images/gapic-showcase:0.16.0 \
    --help
> Root command of gapic-showcase
>
> Usage:
>   gapic-showcase [command]
>
> Available Commands:
>   completion  Emits bash a completion for gapic-showcase
>   echo        This service is used showcase the four main types...
>   help        Help about any command
>   identity    A simple identity service.
>   messaging   A simple messaging service that implements chat...
>   run         Runs the showcase server
>   testing     A service to facilitate running discrete sets of...
>
> Flags:
>   -h, --help      help for gapic-showcase
>   -j, --json      Print JSON output
>   -v, --verbose   Print verbose output
>       --version   version for gapic-showcase
>
> Use "gapic-showcase [command] --help" for more information about a command.
```

### Binary
```sh
$ export GAPIC_SHOWCASE_VERSION=0.16.0
$ export OS=linux
$ export ARCH=amd64
$ curl -L https://github.com/googleapis/gapic-showcase/releases/download/v${GAPIC_SHOWCASE_VERSION}/gapic-showcase-${GAPIC_SHOWCASE_VERSION}-${OS}-${ARCH}.tar.gz | sudo tar -zx --directory /usr/local/bin/
$ gapic-showcase --help
...
```

### Source
```sh
$ go get github.com/googleapis/gapic-showcase/cmd/gapic-showcase
$ PATH=$PATH:`go env GOPATH`/bin
$ gapic-showcase --help
...
```
_* Bear in mind this is not a versioned installation so no versioning guarantees
hold using this installation method._

## Schema
The schema of GAPIC Showcase API can be found in [schema/google/showcase/v1beta1](schema/google/showcase/v1beta1)
Its dependencies can be found in the [googleapis/api-common-protos](https://github.com/googleapis/api-common-protos)
submodule.

## Development Environment
To set up this repository for local development, follow these steps:

1. Install `protoc` from the protobuf [release page](https://github.com/protocolbuffers/protobuf/releases)
or your OS package manager. This API utilizes `proto3_optional`, thus `v3.12.0`
is the minimum supported version of `protoc`.

1. Initialize the `api-common-protos` submodule:
    ```sh
    git submodule update --init --recursive
    ```

1. Install Go
    1. Linux: `sudo apt-get install golang`
    2. Mac, Windows, or other options: Please see the [official set-up docs](https://golang.org/doc/install).

1. Clone this repository.

1. Set up Go protobuf tools:
    ```sh
    go install github.com/golang/protobuf/protoc-gen-go
    go get github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_cli
    go get github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
    ```

1. Export the Go binaries to your environment path.
    ```sh
    PATH=$PATH:`go env GOPATH`/bin
    ```

1. To compile the Showcase binary, as well as associated development utilities in this repository, run the following after you make changes:
    ```sh
    go install ./...
    ```



### Making changes to the protos

If there are any changes to the protobuf files, the generated support code must
be regenerated. This can be done by executing the following command:

    go install ./util/cmd/... && go run ./util/cmd/compile_protos

If successful, you may see changes in the following directories:

* `server/genproto`
* `server/genrest`
* `client/`
* `cmd/gapic-showcase`

Then, update the binaries:

    go install ./...

## Quick Start
This quick start guide will show you how to start the server and make a request to it.

### Step 1. Run the server
Run the showcase server to allow requests to be sent to it. This opens port :7469 to
send and receive requests.

```sh
$ gapic-showcase run
> 2018/09/19 02:13:09 Showcase listening on port: :7469
```

### Step 2. Make a request
Open a new terminal window and make a request to the server.
```sh
$ gapic-showcase \
  identity \ # Service name
  create-user \ # Message name
  --user.display_name Rumble \ # Request fields
  --user.email rumble@goodboi.com
> name:"users/0" display_name:"Rumble" email:"rumble@goodboi.com" create_time:<seconds:1554414332 nanos:494679000 > update_time:<seconds:1554414332 nanos:494679000 >
```
**_Note: You can make requests to this server from your own client but an insecure channel
must be used since the server does not implement auth. Client library generators with
Showcase-based integration tests need to provide the insecure channel to the client library
in the tests._**

Example for Node.js:

```js
const grpc = require('@grpc/grpc-js');
const showcase = require('showcase');
const client = new showcase.EchoClient({ grpc, sslCreds: grpc.credentials.createInsecure() });
```

Example for Go:

```go
conn, err := grpc.Dial("localhost:7469", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
opt := option.WithGRPCConn(conn)
client, err = showcase.NewEchoClient(context.Background(), opt)
if err != nil {
    log.Fatal(err)
}
```

## Released Artifacts
GAPIC Showcase releases three main artifacts, a CLI tool, the gapic-showcase
service protobuf files staged alongside its dependencies, and a protocol buffer
descriptor set compiled from the gapic-showcase service protos.

Check out our [releases](https://github.com/googleapis/gapic-showcase/releases) page to see our released artifacts.

## Versioning
GAPIC Showcase follows semantic versioning. All artifacts that are
released for a certain version are guaranteed to be compatible with one another.

## Supported Go Versions
GAPIC Showcase is supported for go versions 1.11 and later.

## FAQ

### Is this Showcase API publicly served?

This API is not publicly served.

## Disclaimer

This is not an official Google product.

[codecovimg]: https://codecov.io/github/googleapis/gapic-showcase/coverage.svg?branch=master
[codecov]: https://codecov.io/github/googleapis/gapic-showcase?branch=master
[godoc]: https://godoc.org/github.com/googleapis/gapic-showcase/server
[godocimg]: https://godoc.org/github.com/googleapis/gapic-showcase/server?status.svg
[releaselevel]: https://cloud.google.com/terms/launch-stages
[releaselevelimg]: https://img.shields.io/badge/release%20level-beta-red.svg?style&#x3D;flat
