# GAPIC Showcase

[![Release Level][releaselevelimg]][releaselevel]
[![CircleCI][circleimg]][circle]
[![Code Coverage][codecovimg]][codecov]
[![GoDoc][godocimg]][godoc]

> GAPIC (Generated API Client) Showcase: an API to test your gRPC clients.

The GAPIC Showcase API is a representative example of the API client
configurations of Google APIs. GAPIC Showcase provides a tool that allows
a user to spin up a server that accepts gRPC requests. This tool also exposes
an interface to make requests to the server but requests from your own
gRPC clients are accepted.

## Installation
The GAPIC Showcase CLI can be installed using three different mechanisms:
pulling a docker image from [Google Container Registry](https://gcr.io/gapic-images/gapic-showcase),
downloading the compiled binary from our our [releases](https://github.com/googleapis/gapic-showcase/releases)
page, or simply by installing from source using go.

### Docker
```sh
$ docker pull gcr.io/gapic-images/gapic-showcase:0.0.16
$ docker run \
    --rm \
    -p 7469:7469/tcp \
    -p 7469:7469/udp \
    gcr.io/gapic-images/gapic-showcase:0.0.16 \
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
$ export GAPIC_SHOWCASE_VERSION=0.0.16
$ export OS=linux
$ export ARCH=amd64
$ curl -L https://github.com/googleapis/gapic-showcase/releases/download/v${GAPIC_SHOWCASE_VERSION}/gapic-showcase-${GAPIC_SHOWCASE_VERSION}-${OS}-${ARCH} | sudo tar -zx -- --directory /usr/local/bin/
$ gapic-showcase --help
...
```

### Source
```sh
$ go install github.com/googleapis/gapic-showcase/cmd/gapic-showcase
$ gapic-showcase --help
...
```
_* Bear in mind this is not a versioned installation so no versioning guarantees
hold using this installation method._

## Schema
The schema of GAPIC Showcase API can be found in [schema/google/showcase/v1alpha3](schema/google/showcase/v1alpha3)
It's dependencies can be found in the [schema/api-common-protos](schema/api-common-protos)
submodule.

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
_* You can make requests to this server from your own client but an insecure channel
must be used since the server does not implement auth._


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

[circle]: https://circleci.com/gh/googleapis/gapic-showcase
[circleimg]: https://circleci.com/gh/googleapis/gapic-showcase.svg?style=shield
[codecovimg]: https://codecov.io/github/googleapis/gapic-showcase/coverage.svg?branch=master
[codecov]: https://codecov.io/github/googleapis/gapic-showcase?branch=master
[godoc]: https://godoc.org/github.com/googleapis/gapic-showcase/server
[godocimg]: https://godoc.org/github.com/googleapis/gapic-showcase/server?status.svg
[releaselevel]: https://cloud.google.com/terms/launch-stages
[releaselevelimg]: https://img.shields.io/badge/release%20level-alpha-red.svg?style&#x3D;flat
