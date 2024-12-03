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
    mux.HandleFunc("/mcp", s.handleMCP)

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

func (s *server) handleMCP(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    resp, err := s.handler.HandleRequest(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
