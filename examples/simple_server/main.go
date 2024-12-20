package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"

	"github.com/DeanXu2357/mcp-go/pkg/mcp"
)

type exampleHandler struct {
	logger logr.Logger
}

func (h *exampleHandler) HandleRequest(ctx context.Context, req *mcp.Request) (*mcp.Response, error) {
	h.logger.Info("handling request", "action", req.Action)

	// Echo back the request payload
	return &mcp.Response{
		Status: "success",
		Data:   req.Payload,
	}, nil
}

func main() {
	// Create a simple logger
	logger := funcr.New(func(prefix, args string) {
		log.Printf("%s: %s\n", prefix, args)
	}, funcr.Options{})

	// Create server configuration
	config := mcp.Config{
		Port:     8080,
		Host:     "localhost",
		LogLevel: "info",
	}

	// Create handler
	handler := &exampleHandler{
		logger: logger,
	}

	// Create and start server
	server := mcp.NewServer(config, logger, handler)

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
