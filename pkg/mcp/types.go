package mcp

import (
    "context"
    "encoding/json"
)

const (
    // ProtocolVersion is the supported MCP protocol version
    ProtocolVersion = "2024-11-05"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
    JSONRPC string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
    ID      interface{}     `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
    JSONRPC string       `json:"jsonrpc"`
    Result  interface{} `json:"result,omitempty"`
    Error   *RPCError   `json:"error,omitempty"`
    ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC 2.0 error
type RPCError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// InitializeParams represents the parameters for the initialize method
type InitializeParams struct {
    ProtocolVersion string `json:"protocolVersion"`
}

// InitializeResult represents the result of the initialize method
type InitializeResult struct {
    ProtocolVersion string            `json:"protocolVersion"`
    ServerInfo      ServerInfo        `json:"serverInfo"`
    Capabilities    map[string]bool   `json:"capabilities"`
}

// ServerInfo represents information about the server
type ServerInfo struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`
}

// Handler defines the interface for handling MCP requests
type Handler interface {
    // Initialize handles the initialize request
    Initialize(ctx context.Context, params *InitializeParams) (*InitializeResult, error)
    // HandleMethod handles other method calls after initialization
    HandleMethod(ctx context.Context, method string, params json.RawMessage) (interface{}, error)
}
