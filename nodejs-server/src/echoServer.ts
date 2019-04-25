import * as grpc from 'grpc';
import * as Long from 'long';

import * as echoTypes from '../pbjs-genfiles/echo.js';

import showcaseV1Beta1 = echoTypes.google.showcase.v1beta1;
import google = echoTypes.google;
import longrunning = google.longrunning;
import { OperationsServer } from './operationsServer.js';

/**
 * Implements Echo server based on echo.proto.
 */
export class EchoServer {
  private requestCount: number;
  private paginationOutput: Map<string, string[]>;
  private paginationRequest: Map<string, string>;
  private operationsServer: OperationsServer;
  private verbose: boolean;

  constructor(verbose: boolean, operationsServer: OperationsServer) {
    this.requestCount = 0;
    this.paginationOutput = new Map<string, string[]>();
    this.paginationRequest = new Map<string, string>();
    this.operationsServer = operationsServer;
    this.verbose = verbose;
  }

  private log(request: number, ...args: Array<string | {}>) {
    if (this.verbose) {
      console.log(`[EchoServer request #${request}]`, ...args);
    }
  }

  /**
   * Helper function to extract `google.rpc.Status` from some requests that
   * contain this field.
   */
  private static requestToStatus(
    request:
      | showcaseV1Beta1.EchoRequest
      | showcaseV1Beta1.ExpandRequest
      | showcaseV1Beta1.WaitRequest
  ): google.rpc.Status {
    const message =
      request.error && request.error.message ? request.error.message : 'Error';
    const code =
      request.error && request.error.code
        ? request.error.code
        : grpc.status.INVALID_ARGUMENT;
    const details =
      request.error && request.error.details ? request.error.details : [];
    const error = new google.rpc.Status();
    error.code = code;
    error.message = message;
    error.details = details;
    return error;
  }

  /**
   * Helper function to build `grpc.ServiceError` out of `google.rpc.Status`.
   */
  private static requestToError(
    request:
      | showcaseV1Beta1.EchoRequest
      | showcaseV1Beta1.ExpandRequest
      | showcaseV1Beta1.WaitRequest
  ): grpc.ServiceError {
    const status = EchoServer.requestToStatus(request);
    const error = new Error(status.message) as grpc.ServiceError;
    error.code = status.code;
    return error;
  }

  /**
   * Unary call example. Echoes the request.
   */
  echo(
    call: grpc.ServerUnaryCall<showcaseV1Beta1.EchoRequest>,
    callback: grpc.requestCallback<showcaseV1Beta1.EchoResponse>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'echo request:', request);
    if (request.content) {
      const response = new showcaseV1Beta1.EchoResponse();
      response.content = request.content;
      this.log(requestId, 'echo response:', response);
      callback(null, request);
    } else {
      const error = EchoServer.requestToError(request);
      this.log(requestId, 'echo error:', error);
      callback(error);
    }
  }

  /**
   * Server streaming call example. Splits the request string into a list of
   * words returned as a stream.
   */
  expand(
    call: grpc.ServerWriteableStream<showcaseV1Beta1.ExpandRequest>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'expand request:', request);
    this.log(requestId, 'expand writing to stream:');
    if (!request.content) {
      const response = new showcaseV1Beta1.EchoResponse();
      this.log(requestId, 'expand error, response:', response);
      call.write(response);
      call.end();
      return;
    }
    const words = request.content.split(/\s+/);
    for (const word of words) {
      const response = new showcaseV1Beta1.EchoResponse();
      response.content = word;
      this.log(requestId, 'expand response:', response);
      call.write(response);
    }
    this.log(requestId, 'expand completed');
    call.end();
  }

  /**
   * Client streaming call example. Joins the list of words received from a
   * stream into a string.
   */
  collect(
    call: grpc.ServerReadableStream<showcaseV1Beta1.EchoRequest>,
    callback: grpc.requestCallback<showcaseV1Beta1.EchoResponse>
  ): void {
    const requestId = ++this.requestCount;
    const results: string[] = [];
    let error: grpc.ServiceError | undefined;
    this.log(requestId, 'collect reading from stream:');
    call.on('data', (request: showcaseV1Beta1.EchoRequest) => {
      this.log(requestId, 'collect request:', request);
      if (request.content) {
        results.push(request.content);
      } else {
        error = EchoServer.requestToError(request);
      }
    });
    call.on('end', () => {
      if (error) {
        this.log(requestId, 'collect error:', error);
        callback(error);
      } else {
        const response = new showcaseV1Beta1.EchoResponse();
        response.content = results.join(' ');
        this.log(requestId, 'collect response:', response);
        callback(null, response);
      }
    });
  }

  /**
   * Bi-directional streaming example. Sends all the data received from the
   * client stream back to the server stream.
   */
  chat(
    call: grpc.ServerDuplexStream<
      showcaseV1Beta1.EchoRequest,
      showcaseV1Beta1.EchoResponse
    >
  ): void {
    const requestId = ++this.requestCount;
    this.log(requestId, 'chat reading from stream, writing to stream:');
    call.on('data', (request: showcaseV1Beta1.EchoRequest) => {
      this.log(requestId, 'chat request:', request);
      const response = new showcaseV1Beta1.EchoResponse();
      response.content = request.content || '';
      this.log(requestId, 'chat response:', response);
      call.write(response);
    });
    call.on('end', () => {
      this.log(requestId, 'chat completed');
      call.end();
    });
  }

  /**
   * Paged iteration. Splits the input string into words, return resulting words
   * in pages as requested.
   */
  pagedExpand(
    call: grpc.ServerUnaryCall<showcaseV1Beta1.PagedExpandRequest>,
    callback: grpc.requestCallback<showcaseV1Beta1.PagedExpandResponse>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'pagedExpand request:', request);
    if (!request.pageSize) {
      request.pageSize = 1;
    }
    let words: string[] = [];
    if (request.pageToken) {
      let errorMessage: string | undefined;
      if (!this.paginationOutput.has(request.pageToken)) {
        errorMessage = `Bad page token ${request.pageToken}`;
      } else if (
        this.paginationRequest.get(request.pageToken) !== request.content
      ) {
        errorMessage = `Page token does not match the original request`;
      }
      if (errorMessage) {
        const error = new Error(errorMessage) as grpc.ServiceError;
        error.code = grpc.status.INVALID_ARGUMENT;
        callback(error);
        return;
      }
      words = this.paginationOutput.get(request.pageToken) || [];
    } else if (request.content) {
      words = request.content.split(/\s+/);
    }
    const results = words.splice(0, request.pageSize);
    const response = new showcaseV1Beta1.PagedExpandResponse();
    response.responses = results.map(word => {
      const echoResponse = new showcaseV1Beta1.EchoResponse();
      echoResponse.content = word;
      return echoResponse;
    });
    if (words.length > 0) {
      this.paginationOutput.set(requestId.toString(), words);
      this.paginationRequest.set(requestId.toString(), request.content);
      response.nextPageToken = requestId.toString();
    }
    this.log(requestId, 'pagedExpand response:', response);
    callback(null, response);
  }

  /**
   * Converts a pair of (seconds, nanos) (often seen in protobufs) to
   * milliseconds.
   */
  private static toMilliseconds(
    seconds: number | Long | null | undefined,
    nanos: number | null | undefined
  ): number {
    let milliseconds = 0;
    if (Long.isLong(seconds)) {
      milliseconds += (seconds as Long).toNumber() * 1000;
    } else if (seconds) {
      milliseconds += (seconds as number) * 1000;
    }
    if (nanos) {
      milliseconds += nanos / 1000000;
    }
    return milliseconds;
  }

  private static toSecondsNanos(milliseconds: number): [number, number] {
    const seconds = Math.floor(milliseconds / 1000);
    const nanos = (milliseconds - seconds * 1000) * 1000000;
    return [seconds, nanos];
  }

  /**
   * Long running operation.
   */
  wait(
    call: grpc.ServerUnaryCall<showcaseV1Beta1.WaitRequest>,
    callback: grpc.requestCallback<longrunning.Operation>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'wait request:', request);

    // calculate how long to sleep based on the request
    let sleepIntervalMs = 0;
    if (request.endTime) {
      sleepIntervalMs =
        EchoServer.toMilliseconds(
          request.endTime.seconds,
          request.endTime.nanos
        ) - new Date().getTime();
    } else if (request.ttl) {
      sleepIntervalMs = EchoServer.toMilliseconds(
        request.ttl.seconds,
        request.ttl.nanos
      );
    }

    // fill the metadata with our expected sleep end time
    const endTimeMs = new Date().getTime() + sleepIntervalMs;
    const [endSeconds, endNanos] = EchoServer.toSecondsNanos(endTimeMs);
    const metadata = new showcaseV1Beta1.WaitMetadata();
    metadata.endTime = new google.protobuf.Timestamp();
    metadata.endTime.seconds = endSeconds;
    metadata.endTime.nanos = endNanos;

    // encode metadata into google.protobuf.Any
    const metadataAny = new google.protobuf.Any();
    metadataAny.value = showcaseV1Beta1.WaitMetadata.encode(metadata).finish();
    metadataAny.type_url =
      'type.googleapis.com/google.showcase.v1beta1.WaitMetadata';

    // create a new long running operation and return it
    const operation = this.operationsServer.newOperation();
    operation.setMetadata(metadataAny);
    const response = operation.getOperation();
    this.log(requestId, 'wait response:', response);
    callback(null, response);

    // perform actual long running operation
    this.log(requestId, `wait waiting for ${sleepIntervalMs} ms...`);
    setTimeout(() => {
      if (request.success) {
        const response = new showcaseV1Beta1.WaitResponse();
        response.content = request.success.content || '';

        // encode response into google.protobuf.Any
        const responseAny = new google.protobuf.Any();
        responseAny.value = showcaseV1Beta1.WaitResponse.encode(
          response
        ).finish();
        responseAny.type_url =
          'type.googleapis.com/google.showcase.v1beta1.WaitResponse';

        // return result of the long running operation
        this.log(requestId, 'wait result:', response);
        operation.setResponse(responseAny);
      } else if (request.error) {
        const error = EchoServer.requestToStatus(request);
        this.log(requestId, 'wait error:', error);
        operation.setError(error);
      }
    }, sleepIntervalMs);
  }
}
