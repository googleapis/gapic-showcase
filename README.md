# GAPIC Showcase

> An API to showcase Generated API Client features and common API patterns used
by Google.

The GAPIC (Generated API Client) Showcase is an API used to show and describe
common patterns used by Google APIS as well as features used by GAPICs to make
calling Google APIs an enjoyable experience. Each method declared by
the Showcase API serves to show a specific pattern or feature.

This repository also includes a server implementation of the Showcase API that
can be used to verify GAPIC generators.

## Method Types, Patterns and Features

### Unary Methods
> A method that sends a single request and returns a single response.

The *Echo* method is used to show Unary methods. This method
simply returns the content specified by in the API call. If an error is
specified to be returned in the API call, the server will respond with the error
specified.

### Server Side Streaming Methods
> A method that sends a single request but returns multiple responses.

The *Expand* method is used to show server side streaming. This method
splits the given string into words and passes each word back on the stream.
If an error is specified to be returned, the server will respond with the
error code after the all words have been sent on the server.

### Client Side Streaming Methods
> A method that sends multiple requests resulting in a single response.

The *Collect* method is used to show client side streaming. This method
accepts strings on passed on the stream. Upon closing the stream, the
strings passed along the stream will be joined on the ' ' character and
returned to the user. If a client passes an error on the stream, the server
will respond with an error and drop the information passed previously.

### Bidirectional Streaming Methods
> A method that sends multiple requests while receiving multiple responses.

The *Chat* method is used to show bidirectional streaming. This method
simply echos the strings that are passed to it. If an error is passed
on the stream, the server will respond with the error code specified.

### Automatic Timeout Handling
> A GAPIC feature to time out requests that take too long to respond.

The *Timeout* method is used to show how GAPICs handle requests that take too
long to respond. Upon receiving a request, the server
will sleep for the requested amount of time end then return the
response specified in the request. If an error is specified, the server
will respond with the error code specified after sleeping.

### Automatic Retry Handling
> A GAPIC feature to RetryId on errors that are known to be retry-able.

This feature allows methods to automatically retry on error codes that
are known to be safe to retry on.

The *Retry* method is used to test how GAPICs automatically handle retrying.
To test retry features, the user must pass a list of responses to the server
 using the *SetupRetry* method. The *SetupRetry* method will return an ID.
Subsequent requests to *Retry* passing this ID will respond in with
the responses specified in the *SetupRetryRequest*.

### Long Running Operations
> An API pattern for methods that take a long time to complete.

The long running operations pattern is used for requests that can take a
long time complete. Generally the first API request returns an operation which
contain the id of the operation and metadata about the status of the operation.
Using the id, a separate method can be called to get the updated statues
of the operation. Whenever an operation is returned, if the operation is
complete, the response data will be included in the operation.

The *Longrunning* method is used to show the long running operation pattern.
The initial request to this method will specify a time that the operation will
complete. Until the specified completion time, `GetOperation` requests
will return an unfinished operation. Upon reaching the completion time,
subsequent `GetOperation` requests will return an operation containing response
specified in the initial request.

### Pagination
> An API pattern for returning a list of items in pages.

This pagination pattern is used to make it such that the caller can specify the
amount of items to be returned in an API response. The pagination pattern is
most often used on LIST methods.

The *Pagination* method is used to show the pagination pattern.

### Parameter Flattening
> A GAPIC feature to make a method take in the fields of a request object as
parameters.

Parameter flattening is a feature to allow methods to be requested
with the top-level fields of the request object passed as individual
parameters to the request.

The *ParameterFlattening* method is used to show parameter flattening.
This method simply returns the request in order to validate that the
individual parameters are structured in the request object correctly.

### Resource Naming
> A GAPIC feature to add type safety to patterned string fields.

Resource naming is a feature that converts strings with a known pattern
into classes that enforce the well known pattern.

The *ResourceNaming* method is used to show resource naming. This method
simply returns the request in order to validate that the request contained the
correct string when converting from the resource name class.

## FAQ

### Is this Showcase API publicly served?

This API is not publicly served. A server implementation of the Showcase API is
included in this repository.

## Disclaimer

This is not an official Google product.
