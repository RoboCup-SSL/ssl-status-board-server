![travis build](https://travis-ci.org/RoboCup-SSL/ssl-status-board-server.svg?branch=master "travis build status")
[![Go Report Card](https://goreportcard.com/badge/github.com/RoboCup-SSL/ssl-status-board-server?style=flat-square)](https://goreportcard.com/report/github.com/RoboCup-SSL/ssl-status-board-server)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/RoboCup-SSL/ssl-status-board-server)
[![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/RoboCup-SSL/ssl-status-board-server/releases/latest)

# SSL Status Board - Server
The server component of the RoboCup SSL status board implemented in Go

You can find the client component here: https://github.com/RoboCup-SSL/ssl-status-board-client

## Installation

Simply go-get it:
```
go get github.com/RoboCup-SSL/ssl-status-board-server
go get github.com/RoboCup-SSL/ssl-status-board-server/ssl-status-board-proxy
```

## Run

After installation:
```
ssl-status-board-server
```

Or without installation:
```
go run ssl-status-board-server.go
```

## Proxy

If the server is not running in the same network as the referee, the proxy can be used: Enable the proxy in the server
via `server-config.yaml` and run it in the local network. Run the proxy on the remote server. 
The local server will connect to the proxy and the proxy receives client connections and passes data from server to 
client.

You can setup a systemd service to automatically start an instance for each field you want to handle. Example configuration
files are provided in `./systemd`.

## Configuration

The server can be configured with `server-config.yaml`. The proxy requires an `proxy-config.conf`. The location of both files can be
passed via command line. Call the executables with `-h` to get the available arguments.