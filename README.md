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
via `server-config.yaml` and run it in the local network. Run the proxy on the remote server: `ssl-status-board-proxy`.
The local server will connect to the proxy and the proxy receives client connections and passes data from server to 
client.

## Configuration

The server can be configured with `server-config.yaml`. The proxy requires an `proxy-config.conf`. The location of both files can be
passed via command line. Call the executables with `-h` to get the available arguments.