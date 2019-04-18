import * as grpc from 'grpc';
import * as echoTypes from '../pbjs-genfiles/echo.js';
import google = echoTypes.google;
import longrunning = google.longrunning;
import { stringify } from 'querystring';

/**
 * Internal representation of a long running operation
 */
export class Operation {
  private operation: longrunning.Operation;
  private callback: () => void;

  constructor(name: string) {
    this.operation = new longrunning.Operation();
    this.operation.done = false;
    this.operation.name = name;
    this.callback = () => {};
  }

  setCallback(callback: () => void): void {
    this.callback = callback;
  }

  setResponse(response: google.protobuf.Any): void {
    this.operation.response = response;
    this.operation.done = true;
    this.callback();
  }

  setMetadata(metadata: google.protobuf.Any): void {
    this.operation.metadata = metadata;
  }

  setError(error: google.rpc.Status): void {
    this.operation.error = error;
    this.operation.done = true;
    this.callback();
  }

  getOperation(): longrunning.Operation {
    return this.operation;
  }
}

/**
 * Implements long running operations as defined in
 * https://github.com/googleapis/googleapis/blob/master/google/longrunning/operations.proto
 */
export class OperationsServer {
  private operationsCount: number;
  private requestCount: number;
  private operations: Map<string, Operation>;
  private verbose: boolean;

  constructor(verbose: boolean) {
    this.operationsCount = 0;
    this.requestCount = 0;
    this.operations = new Map<string, Operation>();
    this.verbose = verbose;
  }

  private log(request: number, ...args: Array<string | {}>) {
    if (this.verbose) {
      console.log(`[OperationsServer request #${request}]`, ...args);
    }
  }

  newOperation(): Operation {
    const id = ++this.operationsCount;
    const name = `operations/${id}`;
    const operation = new Operation(name);
    this.operations.set(name, operation);
    return operation;
  }

  listOperations(
    call: grpc.ServerUnaryCall<longrunning.ListOperationsRequest>,
    callback: grpc.requestCallback<longrunning.ListOperationsResponse>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'listOperations request:', request);
    const error = new Error(
      'ListOperations is not implemented'
    ) as grpc.ServiceError;
    error.code = grpc.status.UNIMPLEMENTED;
    this.log(requestId, 'listOperations error:', error);
    callback(error);
  }

  getOperation(
    call: grpc.ServerUnaryCall<longrunning.GetOperationRequest>,
    callback: grpc.requestCallback<longrunning.Operation>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'getOperation request:', request);
    const operation = this.operations.get(request.name);
    if (!operation) {
      const error = new Error(
        `Operation ${request.name} does not exist`
      ) as grpc.ServiceError;
      error.code = grpc.status.NOT_FOUND;
      this.log(requestId, 'getOperation error:', error);
      callback(error);
      return;
    }
    const response = operation.getOperation();
    this.log(requestId, 'getOperation response:', response);
    callback(null, response);
  }

  deleteOperation(
    call: grpc.ServerUnaryCall<longrunning.DeleteOperationRequest>,
    callback: grpc.requestCallback<google.protobuf.Empty>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'deleteOperation request:', request);
    const operation = this.operations.get(request.name);
    if (!operation) {
      const error = new Error(
        `Operation ${request.name} does not exist`
      ) as grpc.ServiceError;
      error.code = grpc.status.NOT_FOUND;
      this.log(requestId, 'deleteOperation error:', error);
      callback(error);
      return;
    }
    this.operations.delete(request.name);
    const response = new google.protobuf.Empty();
    this.log(requestId, 'deleteOperation response:', response);
    callback(null, response);
  }

  cancelOperation(
    call: grpc.ServerUnaryCall<longrunning.CancelOperationRequest>,
    callback: grpc.requestCallback<google.protobuf.Empty>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'cancelOperation request:', request);
    const error = new Error(
      'CancelOperation is not implemented'
    ) as grpc.ServiceError;
    error.code = grpc.status.UNIMPLEMENTED;
    this.log(requestId, 'cancelOperation error:', error);
    callback(error);
  }

  waitOperation(
    call: grpc.ServerUnaryCall<longrunning.WaitOperationRequest>,
    callback: grpc.requestCallback<longrunning.Operation>
  ): void {
    const requestId = ++this.requestCount;
    const request = call.request;
    this.log(requestId, 'waitOperation request:', request);
    const operation = this.operations.get(request.name);
    if (!operation) {
      const error = new Error(
        `Operation ${request.name} does not exist`
      ) as grpc.ServiceError;
      error.code = grpc.status.NOT_FOUND;
      this.log(requestId, 'waitOperation error:', error);
      callback(error);
      return;
    }
    const response = operation.getOperation();
    if (response.done) {
      this.log(requestId, 'getOperation response:', response);
      callback(null, response);
      return;
    }
    this.log(requestId, 'getOperation waiting...');
    operation.setCallback(() => {
      const response = operation.getOperation();
      this.log(requestId, 'getOperation response:', response);
      callback(null, response);
    });
  }
}
