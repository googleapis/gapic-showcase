import * as protoLoader from '@grpc/proto-loader';
import * as grpc from 'grpc';
import * as path from 'path';
import * as yargs from 'yargs';
import {Argv} from 'yargs';
import {EchoServer} from './echoServer';
import {OperationsServer} from './operationsServer';

const protoRoot = path.join(__dirname, '..', '..', '..', 'schema');
const commonRoot =
    path.join(__dirname, '..', '..', 'node_modules', 'google-proto-files');
const protoPath = path.join('google', 'showcase', 'v1beta1', 'echo.proto');

function loadProtos() {
  const options = {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
    includeDirs: [protoRoot, commonRoot]
  };
  const packageDefinition = protoLoader.loadSync(protoPath, options);
  return packageDefinition;
}

function createServer(verbose: boolean) {
  const server = new grpc.Server();
  const packageDefinition = loadProtos();
  const descriptor = grpc.loadPackageDefinition(packageDefinition);
  const operationsServer = new OperationsServer(verbose);
  const echoServer = new EchoServer(verbose, operationsServer);
  server.addService(
      // @ts-ignore unknown types
      descriptor.google.showcase.v1beta1.Echo.service, echoServer);
  server.addService(
      // @ts-ignore unknown types
      descriptor.google.longrunning.Operations.service, operationsServer);
  return server;
}

function main() {
  const argv = yargs.argv;

  if (argv['help'] || argv['usage']) {
    console.log('Command line options: ');
    console.log('--bind: address:port to bind, e.g. 0.0.0.0:7469');
    console.log('--verbose: be verbose');
    process.exit(1);
  }

  const bindAddress = (argv['bind'] || '0.0.0.0:7469') as string;
  const verbose = argv['verbose'] as boolean;

  const server = createServer(verbose);
  const port =
      server.bind(bindAddress, grpc.ServerCredentials.createInsecure());
  if (port <= 0) {
    console.log(`Failed to bind on ${bindAddress}, exiting.`);
    process.exit(1);
  }
  console.log(`Server is listening on port ${port}.`);
  server.start();
}

main();
