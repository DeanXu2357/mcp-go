package mcp

import (
	"context"
)

// Config represents the MCP server configuration
type Config struct {
	Port     int
	Host     string
	LogLevel string
}

// Server represents an MCP server instance
type Server interface {
	// Start starts the MCP server
	Start(ctx context.Context) error
	// Stop stops the MCP server
	Stop(ctx context.Context) error
}

// Handler defines the interface for handling MCP requests
type Handler interface {
	// HandleRequest handles an incoming MCP request
	HandleRequest(ctx context.Context, req *Request) (*Response, error)
}

// Request represents an MCP request
type Request struct {
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

// Response represents an MCP response
type Response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Error  string                 `json:"error,omitempty"`
}
