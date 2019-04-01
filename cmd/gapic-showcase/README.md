# gapic-showcase
This directory contains the command line interface (CLI) used to run the
Showcase API server as well as make requests to the Showcase API.

## Installation
```sh
$ go install github.com/googleapis/gapic-showcase/cmd/gapic-showcase
```

## Usage
```sh
$ gapic-showcase --help

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

### Running the server
```sh
$ gapic-showcase run --help

> Runs the showcase server
>
> Usage:
>   gapic-showcase run [flags]
>
> Flags:
>   -h, --help          help for run
>   -p, --port string   The port that showcase will be served on. (default ":7469")
>
> Global Flags:
>   -j, --json      Print JSON output
>   -v, --verbose   Print verbose output


$ gapic-showcase run --port 1234

> 2019/04/01 12:36:35 Showcase listening on port: :1234
```

### Making a request
A request can also be made to the Showcase API using this CLI. The command to make a request
is done by using a service's subcommand, the method's subcommand and passing the request values
as flags as shown below.
```
$ gapic-showcase {service_name} {method_name} --{request_field_name} {value}
```

#### Example
```sh
$ gapic-showcase identity --help

> A simple identity service.
>
> Usage:
>   gapic-showcase identity [command]
>
> Available Commands:
>   create-user Creates a user.
>   delete-user Deletes a user, their profile, and all of their...
>   get-user    Retrieves the User with the given uri.
>   list-users  Lists all users.
>   update-user Updates a user.
>
> Flags:
>       --address string   Set API address used by client. Or use GAPIC-SHOWCASE_IDENTITY_ADDRESS.
>       --api_key string   Set API Key used by the client. Or use GAPIC-SHOWCASE_IDENTITY_API_KEY.
>   -h, --help             help for identity
>       --insecure         Make insecure client connection. Or use GAPIC-SHOWCASE_IDENTITY_INSECURE. Must be used with "address" option
>       --token string     Set Bearer token used by the client. Or use GAPIC-SHOWCASE_IDENTITY_TOKEN.
>
> Global Flags:
>   -j, --json      Print JSON output
>   -v, --verbose   Print verbose output
>
> Use "gapic-showcase identity [command] --help" for more information about a command.


$ gapic-showcase identity create-user --help

> Creates a user.
>
> Usage:
>   gapic-showcase identity create-user [flags]
>
> Flags:
>       --from_file string               Absolute path to JSON file containing request payload
>   -h, --help                           help for create-user
>       --user.create_time.nanos int32
>       --user.create_time.seconds int
>       --user.display_name string
>       --user.email string
>       --user.name string
>       --user.update_time.nanos int32
>       --user.update_time.seconds int
>
> Global Flags:
>       --address string   Set API address used by client. Or use GAPIC-SHOWCASE_IDENTITY_ADDRESS.
>       --api_key string   Set API Key used by the client. Or use GAPIC-SHOWCASE_IDENTITY_API_KEY.
>       --insecure         Make insecure client connection. Or use GAPIC-SHOWCASE_IDENTITY_INSECURE. Must be used with "address" option
>   -j, --json             Print JSON output
>       --token string     Set Bearer token used by the client. Or use GAPIC-SHOWCASE_IDENTITY_TOKEN.
>   -v, --verbose          Print verbose output


$ gapic-showcase identity create-user --user.display_name Rumble --user.email rumble@goodboi.com

> name:"users/0" display_name:"Rumble" email:"rumble@goodboi.com" create_time:<seconds:1554144706 nanos:304080000 > update_time:<seconds:1554144706 nanos:304080000 >
```