# gURL
[![Build Status](https://circleci.com/gh/fullstorydev/grpcurl/tree/master.svg?style=svg)](https://circleci.com/gh/fullstorydev/grpcurl/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fullstorydev/grpcurl)](https://goreportcard.com/report/github.com/fullstorydev/grpcurl)
[![Snap Release Status](https://snapcraft.io/grpcurl/badge.svg)](https://snapcraft.io/grpcurl)

`gurl` is a unified command-line tool that supports both gRPC and HTTP/HTTPS requests.

## Overview

`gurl` combines the functionality of `grpcurl` (for gRPC servers) and `curl` (for HTTP/HTTPS endpoints) into a single tool:

- **gRPC mode** (default): Interact with gRPC servers using JSON/text input
- **HTTP mode** (`--curl`): Acts as a standard HTTP client like `curl` for REST APIs and web services

### gRPC Mode (Default)

In standard mode, `gurl` is essentially `grpcurl` - a tool for interacting with gRPC servers. It's
basically `curl` for gRPC servers.

The main purpose is to invoke RPC methods on a gRPC server from the
command-line. gRPC servers use a binary encoding on the wire
([protocol buffers](https://developers.google.com/protocol-buffers/), or "protobufs"
for short). So they are basically impossible to interact with using regular `curl`
(and older versions of `curl` that do not support HTTP/2 are of course non-starters).
This program accepts messages using JSON encoding, which is much more friendly for both
humans and scripts.

With this tool you can also browse the schema for gRPC services, either by querying
a server that supports [server reflection](https://github.com/grpc/grpc/blob/master/src/proto/grpc/reflection/v1/reflection.proto),
by reading proto source files, or by loading in compiled "protoset" files (files that contain
encoded file [descriptor protos](https://github.com/google/protobuf/blob/master/src/google/protobuf/descriptor.proto)).
In fact, the way the tool transforms JSON request data into a binary encoded protobuf
is using that very same schema. So, if the server you interact with does not support
reflection, you will either need the proto source files that define the service or need
protoset files that `gurl` can use.

This repo also provides a library package, `github.com/fullstorydev/grpcurl`, that has
functions for simplifying the construction of other command-line tools that dynamically
invoke gRPC endpoints. This code is a great example of how to use the various packages of
the [protoreflect](https://godoc.org/github.com/jhump/protoreflect) library, and shows
off what they can do.

See also the [`grpcurl` talk at GopherCon 2018](https://www.youtube.com/watch?v=dDr-8kbMnaw).

### HTTP Mode (--curl)

When `--curl` is enabled, `gurl` acts as a standard HTTP/HTTPS client, making it compatible
with regular web services and REST APIs. This mode supports common curl options and behavior.

## Features

### gRPC Mode Features
`gurl` supports all kinds of RPC methods, including streaming methods. You can even
operate bi-directional streaming methods interactively by running `gurl` from an
interactive terminal and using stdin as the request body!

`gurl` supports both secure/TLS servers _and_ plain-text servers (i.e. no TLS) and has
numerous options for TLS configuration. It also supports mutual TLS, where the client is
required to present a client certificate.

As mentioned above, `gurl` works seamlessly if the server supports the reflection
service. If not, you can supply the `.proto` source files or you can supply protoset
files (containing compiled descriptors, produced by `protoc`) to `gurl`.

### cURL Mode Features
- Standard HTTP methods: GET, POST, PUT, DELETE, etc.
- Custom headers with `-H` or `--header`
- Request body data with `-d` or `--data`
- TLS support with certificate validation
- Client certificates for mutual TLS
- Unix domain socket support
- Verbose output with `-v`
- Output to file with `-o`
- And more curl-compatible options

## Installation

### From Source
If you already have the [Go SDK](https://golang.org/doc/install) installed, you can build `gurl` from source:
```shell
# Clone this repository
git clone https://github.com/worthies/gURL.git
cd gURL

# Build the binary
go build -o gurl ./cmd/gurl

# Optionally install to $GOPATH/bin
go install ./cmd/gurl
```

This installs the command into the `bin` sub-folder of wherever your `$GOPATH`
environment variable points. (If you have no `GOPATH` environment variable set,
the default install location is `$HOME/go/bin`). If this directory is already in
your `$PATH`, then you should be good to go.

## Usage

### Basic Usage

The usage doc for the tool explains the numerous options:
```shell
gurl -help
```

### gRPC Mode (Default)

In the sections below, you will find numerous examples demonstrating how to use
`gurl` in gRPC mode (default behavior).

#### Invoking RPCs
Invoking an RPC on a trusted server (e.g. TLS without self-signed key or custom CA)
that requires no client certs and supports server reflection is the simplest thing to
do with `gurl`. This minimal invocation sends an empty request body:
```shell
gurl grpc.server.com:443 my.custom.server.Service/Method

# no TLS
gurl -plaintext grpc.server.com:80 my.custom.server.Service/Method
```

To send a non-empty request, use the `-d` argument. Note that all arguments must come
*before* the server address and method name:
```shell
gurl -d '{"id": 1234, "tags": ["foo","bar"]}' \
    grpc.server.com:443 my.custom.server.Service/Method
```

As can be seen in the example, the supplied body must be in JSON format. The body will
be parsed and then transmitted to the server in the protobuf binary format.

If you want to include `gurl` in a command pipeline, such as when using `jq` to
create a request body, you can use `-d @`, which tells `gurl` to read the actual
request body from stdin:
```shell
gurl -d @ grpc.server.com:443 my.custom.server.Service/Method <<EOM
{
  "id": 1234,
  "tags": [
    "foor",
    "bar"
  ]
}
EOM
```
### Adding Headers/Metadata to Request
Adding of headers / metadata to a rpc request is possible via the `-H name:value` command line option. Multiple headers can be added in a similar fashion.
Example :
```shell
gurl -H header1:value1 -H header2:value2 -d '{"id": 1234, "tags": ["foo","bar"]}' grpc.server.com:443 my.custom.server.Service/Method
```
For more usage guide, check out the help docs via `gurl -help`

### cURL-compatible mode

`gurl` can operate in a cURL-compatible mode where common cURL-style
options are accepted and mapped to `gurl`'s flags. Enable this mode with
`--curl`. The following mappings are supported:

#### Header and Data Options
- `-H, --header <name: value>` → `-H "name: value"` (pass custom headers)
- `-d, --data <data>` → `-d <data>` (request body data)
- `--data-raw, --data-binary, --data-ascii` → `-d <data>` (all treated as request data)

#### TLS/Security Options
- `-k, --insecure` → `-insecure` (skip TLS verification)
- `--cacert <file>` → `-cacert <file>` (CA certificate)
- `-E, --cert <cert[:password]>` → `-cert <cert>` (client certificate)
- `--key <file>` → `-key <file>` (private key)

#### Connection Options
- `--connect-timeout <seconds>` → `-connect-timeout <seconds>` (connection timeout)
- `-m, --max-time <seconds>` → `-max-time <seconds>` (maximum request time)
- `--keepalive-time <seconds>` → `-keepalive-time <seconds>` (keepalive interval)
- `--unix-socket <path>` → sets `-unix` and uses the socket path as address
- `--abstract-unix-socket <path>` → same as `--unix-socket`

#### General Options
- `-A, --user-agent <name>` → `-user-agent <name>` (set User-Agent header)
- `-u, --user <user:password>` → adds `Authorization: Basic` header
- `-v, --verbose` → `-v` (verbose output)
- `-h, --help` → `-help` (show help)
- `-V, --version` → `-version` (show version)

#### URL Handling
- `--url <url>` → extracts host:port from URL
- URLs starting with `http://` or `https://` → automatically extracts host:port

#### Combined Short Flags
Short flags can be combined: `-kv` is equivalent to `-k -v` (insecure + verbose)

#### Unsupported Options
The following cURL options are not applicable to gRPC and are silently ignored:
- HTTP-specific: `--compressed`, `--location`, `--cookie`, `--form`, `--get`, `--head`
- Output options: `--output`, `--remote-name`, `--dump-header`
- Proxy options: `--proxy`, `--noproxy`, `--proxy-user`
- Protocol versions: `--http1.0`, `--http1.1`, `--http2`

Examples:

```shell
# Basic usage with headers and data
gurl --curl --header "authorization: Bearer TOKEN" \
    --data '{"id":1}' grpc.server.com:443 my.Service/Method

# Using short flags (combined)
gurl --curl -kv -d '{"name":"test"}' localhost:50051 my.Service/Method

# Unix socket
gurl --curl --unix-socket /tmp/grpc.sock my.Service/Method

# With client certificates
gurl --curl --cert client.crt --key client.key \
    --cacert ca.crt grpc.server.com:443 my.Service/Method

# URL-style addressing
gurl --curl -d '{"id":1}' https://api.example.com/my.Service/Method

# With timeout and keepalive
gurl --curl --connect-timeout 5 --max-time 30 \
    --keepalive-time 10 server:443 my.Service/Method
```### Listing Services
To list all services exposed by a server, use the "list" verb. When using `.proto` source
or protoset files instead of server reflection, this lists all services defined in the
source or protoset files.
```shell
# Server supports reflection
gurl localhost:8787 list

# Using compiled protoset files
gurl -protoset my-protos.bin list

# Using proto sources
gurl -import-path ../protos -proto my-stuff.proto list

# Export proto files (use -proto-out-dir to specify the output directory)
gurl -plaintext -proto-out-dir "out_protos" "localhost:8787" describe my.custom.server.Service

# Export protoset file (use -protoset-out to specify the output file)
gurl -plaintext -protoset-out "out.protoset" "localhost:8787" describe my.custom.server.Service

```

The "list" verb also lets you see all methods in a particular service:
```shell
gurl localhost:8787 list my.custom.server.Service
```

### Describing Elements
The "describe" verb will print the type of any symbol that the server knows about
or that is found in a given protoset file. It also prints a description of that
symbol, in the form of snippets of proto source. It won't necessarily be the
original source that defined the element, but it will be equivalent.

```shell
# Server supports reflection
gurl localhost:8787 describe my.custom.server.Service.MethodOne

# Using compiled protoset files
gurl -protoset my-protos.bin describe my.custom.server.Service.MethodOne

# Using proto sources
gurl -import-path ../protos -proto my-stuff.proto describe my.custom.server.Service.MethodOne
```

## Descriptor Sources
The `gurl` tool can operate on a variety of sources for descriptors. The descriptors
are required, in order for `gurl` to understand the RPC schema, translate inputs
into the protobuf binary format as well as translate responses from the binary format
into text. The sections below document the supported sources and what command-line flags
are needed to use them.

### Server Reflection

Without any additional command-line flags, `gurl` will try to use [server reflection](https://github.com/grpc/grpc/blob/master/src/proto/grpc/reflection/v1/reflection.proto).

Examples for how to set up server reflection can be found [here](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md#known-implementations).

When using reflection, the server address (host:port or path to Unix socket) is required
even for "list" and "describe" operations, so that `gurl` can connect to the server
and ask it for its descriptors.

### Proto Source Files
To use `gurl` on servers that do not support reflection, you can use `.proto` source
files.

In addition to using `-proto` flags to point `gurl` at the relevant proto source file(s),
you may also need to supply `-import-path` flags to tell `gurl` the folders from which
dependencies can be imported.

Just like when compiling with `protoc`, you do *not* need to provide an import path for the
location of the standard protos included with `protoc` (which contain various "well-known
types" with a package definition of `google.protobuf`). These files are "known" by `gurl`
as a snapshot of their descriptors is built into the `gurl` binary.

When using proto sources, you can omit the server address (host:port or path to Unix socket)
when using the "list" and "describe" operations since they only need to consult the proto
source files.

### Protoset Files
You can also use compiled protoset files with `gurl`. If you are scripting `gurl` and
need to re-use the same proto sources for many invocations, you will see better performance
by using protoset files (since it skips the parsing and compilation steps with each
invocation).

Protoset files contain binary encoded `google.protobuf.FileDescriptorSet` protos. To create
a protoset file, invoke `protoc` with the `*.proto` files that define the service:
```shell
protoc --proto_path=. \
    --descriptor_set_out=myservice.protoset \
    --include_imports \
    my/custom/server/service.proto
```

The `--descriptor_set_out` argument is what tells `protoc` to produce a protoset,
and the `--include_imports` argument is necessary for the protoset to contain
everything that `gurl` needs to process and understand the schema.

When using protosets, you can omit the server address (host:port or path to Unix socket)
when using the "list" and "describe" operations since they only need to consult the
protoset files.

