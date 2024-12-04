# MCP Go SDK

A Go SDK for implementing Model Context Protocol (MCP) servers.

## Requirements

- Go 1.23 or later

## Installation

```bash
go get github.com/DeanXu2357/mcp-go
```

## Quick Start

Here's a simple example of how to create an MCP server:

```go
package main

import (
    "github.com/DeanXu2357/mcp-go/pkg/mcp"
)

func main() {
    // Create server configuration
    config := mcp.Config{
        Port:     8080,
        Host:     "localhost",
        LogLevel: "info",
    }

    // Create and start server
    server := mcp.NewServer(config, logger, handler)
    server.Start(context.Background())
}
```

For a complete example, see the [examples/simple_server](examples/simple_server) directory.

## Features

- Easy-to-use API for implementing MCP servers
- Built-in logging using go-logr
- Graceful shutdown support
- Extensible handler interface

## License

MIT License
