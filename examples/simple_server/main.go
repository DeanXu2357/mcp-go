package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/go-logr/logr"
    "github.com/go-logr/logr/funcr"
    "github.com/yourusername/mcp_go_sdk/pkg/mcp"
)

type exampleHandler struct {
    logger logr.Logger
    initialized bool
}

func (h *exampleHandler) Initialize(ctx context.Context, params *mcp.InitializeParams) (*mcp.InitializeResult, error) {
    h.logger.Info("initializing server", "protocolVersion", params.ProtocolVersion)
    
    // Check protocol version compatibility
    if params.ProtocolVersion != mcp.ProtocolVersion {
        h.logger.Info("protocol version mismatch", 
            "client", params.ProtocolVersion,
            "server", mcp.ProtocolVersion)
    }

    h.initialized = true
    
    return &mcp.InitializeResult{
        ProtocolVersion: mcp.ProtocolVersion,
        ServerInfo: mcp.ServerInfo{
            Name:        "example-server",
            Version:     "0.1.0",
            Description: "An example MCP server",
        },
        Capabilities: map[string]bool{
            "echo": true,
        },
    }, nil
}

func (h *exampleHandler) HandleMethod(ctx context.Context, method string, params json.RawMessage) (interface{}, error) {
    // Ensure server is initialized
    if !h.initialized {
        return nil, &mcp.RPCError{
            Code:    -32002,
            Message: "Server not initialized",
        }
    }

    h.logger.Info("handling method", "method", method)
    
    switch method {
    case "echo":
        var data map[string]interface{}
        if err := json.Unmarshal(params, &data); err != nil {
            return nil, err
        }
        return data, nil
    default:
        return nil, &mcp.RPCError{
            Code:    -32601,
            Message: "Method not found",
        }
    }
}

func main() {
    // Create a simple logger
    logger := funcr.New(func(prefix, args string) {
        log.Printf("%s: %s\n", prefix, args)
    }, funcr.Options{})

    // Create server configuration
    serverConfig := mcp.ServerConfig{
        Name: "example-server",
        Type: "example",
        Host: mcp.DefaultHost,
        Port: mcp.DefaultPort,
        Capabilities: map[string]bool{
            "echo": true,
        },
    }

    // Register server with Claude Desktop
    if err := mcp.RegisterServer(serverConfig); err != nil {
        logger.Error(err, "failed to register server")
        os.Exit(1)
    }

    // Create handler
    handler := &exampleHandler{
        logger: logger,
    }

    // Create and start server
    server := mcp.NewServer(mcp.Config{
        Port:     serverConfig.Port,
        Host:     serverConfig.Host,
        LogLevel: "info",
    }, logger, handler)

    // Setup context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle shutdown gracefully
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
        <-sigCh
        logger.Info("shutting down server...")
        cancel()
    }()

    // Start server
    if err := server.Start(ctx); err != nil {
        logger.Error(err, "server error")
        os.Exit(1)
    }
}
