# Feature Testing Server

> A server to test generated clients.

This is a server used to test the features of clients generated
by the gapic-generators in [googleapis]. Each method declared by
the Feature Testing API serves to test one feature of a client.
Below is listed the features that can be tested using this API
alongside test cases that should be tested to verify the client
feature in question.

[googleapis](https://github.com/googleapis)

## Features

### Unary

The *Echo* method is used to test Unary calls and responses. This
method simply returns the content specified by in the API call.
If an error is specified to be returned in the API call, the
server will respond with the error code specified.

#### Test Cases

- Single string request
- Server error

### Server Side Streaming

The *Expand* method is used to test server side streaming. This method
splits the given string into words and passes each word back on the stream.
If an error is specified to be returned, the server will respond with the
error code after the all words have been sent on the server.

#### Test Cases

- Simple string request
- Error on first request
- Error after successful stream responses

### Client Side Streaming

The *Collect* method is used to test client side streaming. This method
accepts strings on passed on the stream. Upon closing the stream, the
strings passed along the stream will be joined on the ' ' character and
returned to the user. If a client passes an error on the stream, the server
will respond with an error and drop the information passed previously.

#### Test Cases

- Passing two strings on the stream
- Error on first request
- Error after strings have been passed on the stream

### Bidirectional Streaming

The *Chat* method is used to test bidirectional streaming. This method
simply echos the strings that are passed to it. If an error is passed
on the stream, the server will respond with the error code specified.

#### Test Cases

- Simple request
- Error on first request
- Error after strings have been passed on the stream.

### Timeout

This feature allows the user to specify an amount of time
that a request should wait on a request.

The *TimeoutTest* method is used to test how clients handle methods
that take a long time to return. Upon receiving a request, the server
will sleep for the requested amount of time end then return the
response specified in the request. If an error is specified, the server
will respond with the error code specified after sleeping.

#### Test Cases

- Responding immediately with success
- Responding after the client should timeout with success
- Responding immediately with an error
- Responding after the client should timeout with an error

### Automatic Retry

This feature allows methods to automatically retry on error codes that
are known to be safe to retry on.

The *RetryTest* method is used to test how clients retry. To test retry
features, the user must pass a list of responses to the server using the
*SetupRetryTest* method. The *SetupRetryTest* method will return an ID.
Subsequent requests to *RetryTest* passing this ID will respond in with
the responses specified in the *SetupRetryTestRequest*.

#### Test Cases

- Successful at first call
- Multiple retry-able errors terminating in success
- Multiple retry-able errors terminating with an non retry-able error

### Longrunning Operations

The longrunning operations pattern is used for requests that can take a
long time. Generally the first request returns an operation which contain
the id of the operation and metadata about the status of the operation.
Using the id, a separate method can be called to get the updated statues
of the operation. Whenever an operation is returned, if the operation is
complete, the response data will be included in the operation.

The *LongrunningTest* method is used to test the how a client handles
longrunning operations. The request will specify a time that the
longrunning operation will complete. Until that completion time, the
requests for the operation will indicate that the operation is not done.
Upon reaching the completion time, subsequent requests for the operation
will return an operation containing response specified in the
initial request.

#### Test Cases

- Operation success on first call
- Operation success after small amount of time.
- Operation success after time to live limit is exceeded for a channel
- Operation error on first call
- Operation error after small amount of time.
- Operation metadata handling


### Pagination

Pagination is a pattern of API call in which the user specifies the
amount of items to be returned in a list as well as the location
in the list to return.

The *PaginationTest* method is used to test client features built
around the pagination pattern.

#### Test Cases

- Successful pagination up to response limit.
- API Call returning fewer items than page size
- API Call returning more items than page size.

### Parameter Flattening

Parameter flattening is a client feature to allow methods to be requested
with the top-level fields of the request object passed as individual
parameters to the request.

The *ParameterFlatteningTest* method is used to test how parameter flattening.
This method simply returns the request in order to validate that the
individual parameters are structured in the request object correctly.

#### Test Cases

- Single Field Flattening
- Repeated Field Flattening
- Nested Object Flattening

### Resource Naming

Resource naming is a feature that converts strings with a well known pattern
into classes that enforce the well known pattern.

The *ResourceNamingTest* method is used to test resource naming. This method
simply returns the request in order to validate that the request contained the
correct string when converting from the resource name class.

### Test Cases

- Field which has a single pattern
- Field which can have multiple patterns

## Disclaimer

This is not an official Google product.
