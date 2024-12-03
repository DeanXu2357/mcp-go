package mcp

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/go-logr/logr"
)

type server struct {
    config  Config
    logger  logr.Logger
    handler Handler
    srv     *http.Server
}

// NewServer creates a new MCP server instance
func NewServer(config Config, logger logr.Logger, handler Handler) Server {
    return &server{
        config:  config,
        logger:  logger,
        handler: handler,
    }
}

func (s *server) Start(ctx context.Context) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/", s.handleJSONRPC)

    s.srv = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
        Handler: mux,
    }

    s.logger.Info("starting MCP server", "addr", s.srv.Addr)
    
    go func() {
        <-ctx.Done()
        s.Stop(context.Background())
    }()

    if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
        return fmt.Errorf("failed to start server: %w", err)
    }

    return nil
}

func (s *server) Stop(ctx context.Context) error {
    s.logger.Info("stopping MCP server")
    return s.srv.Shutdown(ctx)
}

func (s *server) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req JSONRPCRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSONRPCError(w, &JSONRPCResponse{
            JSONRPC: "2.0",
            Error: &RPCError{
                Code: -32700,
                Message: "Parse error",
            },
            ID: nil,
        })
        return
    }

    // Validate JSON-RPC version
    if req.JSONRPC != "2.0" {
        writeJSONRPCError(w, &JSONRPCResponse{
            JSONRPC: "2.0",
            Error: &RPCError{
                Code: -32600,
                Message: "Invalid Request",
            },
            ID: req.ID,
        })
        return
    }

    // Handle the method
    result, err := s.handler.HandleMethod(r.Context(), req.Method, req.Params)
    if err != nil {
        writeJSONRPCError(w, &JSONRPCResponse{
            JSONRPC: "2.0",
            Error: &RPCError{
                Code: -32603,
                Message: err.Error(),
            },
            ID: req.ID,
        })
        return
    }

    // Send response
    response := &JSONRPCResponse{
        JSONRPC: "2.0",
        Result:  result,
        ID:      req.ID,
    }

    writeJSONRPCResponse(w, response)
}

func writeJSONRPCError(w http.ResponseWriter, response *JSONRPCResponse) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func writeJSONRPCResponse(w http.ResponseWriter, response *JSONRPCResponse) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
