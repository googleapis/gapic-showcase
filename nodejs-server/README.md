## Node.js gRPC server

This is Node.js gRPC server that can serve RPCs defined in 
[echo.proto](https://github.com/googleapis/gapic-showcase/blob/master/schema/google/showcase/v1beta1/echo.proto) 
(including long running operations). 

This is an experimental showcase server implementation. It will eventually
be upgraded to support all services according to the spec, but for now please
use 
[the Go implementation](https://github.com/googleapis/gapic-showcase/tree/master/server)
(the one shipped in the Docker image) if you need anything more than Echo service.

### Usage

Start server:
```sh
$ npm install
$ node build/src/index.js --verbose    # will start on port 7469
```

Play with client:
```sh
$ docker run --rm --network host gcr.io/gapic-images/gapic-showcase:0.2.1 \
  echo --address host.docker.internal:7469 echo --response content --response.content okay

$ docker run --rm --network host gcr.io/gapic-images/gapic-showcase:0.2.1 \
  echo --address host.docker.internal:7469 wait --end ttl --end.ttl.seconds 5 \
  --follow --response success --response.success.content okay
```
