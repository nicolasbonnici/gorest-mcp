package mcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest-mcp/middleware"
)

// ConnectionPool tracks active SSE connections per user
type ConnectionPool struct {
	mu          sync.RWMutex
	connections map[string]int // user_id -> connection count
}

var connectionPool = &ConnectionPool{
	connections: make(map[string]int),
}

// AddConnection adds a new connection for a user
func (cp *ConnectionPool) AddConnection(userID string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	count := cp.connections[userID]
	if count >= 3 { // Max connections from config
		return ErrMaxConnectionsExceeded
	}

	cp.connections[userID]++
	return nil
}

// RemoveConnection removes a connection for a user
func (cp *ConnectionPool) RemoveConnection(userID string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if count := cp.connections[userID]; count > 0 {
		cp.connections[userID]--
		if cp.connections[userID] == 0 {
			delete(cp.connections, userID)
		}
	}
}

func HandleSSE(c *fiber.Ctx, mcpServer *MCPServer, config *Config, log *slog.Logger) error {
	// Validate JWT and extract user context
	userID, userEmail, err := middleware.ValidateJWT(c, log)
	if err != nil {
		log.Error("JWT validation failed", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid or missing JWT token",
		})
	}

	// Check connection limits
	if err := connectionPool.AddConnection(userID); err != nil {
		log.Warn("Max connections exceeded", "user_id", userID, "email", userEmail)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Maximum concurrent connections exceeded",
		})
	}
	defer connectionPool.RemoveConnection(userID)

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Create context with user info
	ctx := context.WithValue(c.Context(), ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, ContextKeyUserEmail, userEmail)

	log.Info("SSE connection established", "user_id", userID, "email", userEmail)

	// Send initial connected event
	if _, err := fmt.Fprintf(c, "event: connected\ndata: {\"status\":\"ready\"}\n\n"); err != nil {
		log.Error("Failed to send connected event", "error", err)
		return err
	}

	// Handle SSE transport using mark3labs/mcp-go
	return handleMCPSSETransport(ctx, c, mcpServer.GetServer(), config, log, userID, userEmail)
}

// handleMCPSSETransport manages the MCP SSE transport lifecycle
func handleMCPSSETransport(
	ctx context.Context,
	c *fiber.Ctx,
	mcpServer *server.MCPServer,
	config *Config,
	log *slog.Logger,
	userID, userEmail string,
) error {
	// Create heartbeat ticker
	heartbeat := time.NewTicker(config.GetHeartbeatInterval())
	defer heartbeat.Stop()

	// Create idle timeout
	idleTimeout := time.NewTimer(config.GetConnectionTimeout())
	defer idleTimeout.Stop()

	// Connection management
	done := make(chan struct{})
	defer close(done)

	// Heartbeat loop
	go func() {
		for {
			select {
			case <-heartbeat.C:
				// Send heartbeat event
				if _, err := fmt.Fprintf(c, "event: heartbeat\ndata: {\"timestamp\":\"%s\"}\n\n",
					time.Now().Format(time.RFC3339)); err != nil {
					log.Error("Failed to send heartbeat", "error", err, "user_id", userID)
					return
				}

				// Reset idle timeout
				if !idleTimeout.Stop() {
					select {
					case <-idleTimeout.C:
					default:
					}
				}
				idleTimeout.Reset(config.GetConnectionTimeout())

			case <-done:
				return
			}
		}
	}()

	// Note: For v0.1.0, we're providing a basic SSE connection with heartbeat.
	// Full MCP request/response handling via SSE will be implemented by integrating
	// mark3labs/mcp-go SSE transport in future versions.
	// Currently, clients can maintain connection and receive heartbeats.

	log.Info("SSE connection ready for MCP communication", "user_id", userID, "email", userEmail)

	// Keep connection alive until timeout or client disconnect
	select {
	case <-idleTimeout.C:
		log.Info("SSE connection idle timeout", "user_id", userID)
		if _, err := fmt.Fprintf(c, "event: timeout\ndata: {\"reason\":\"idle_timeout\"}\n\n"); err != nil {
			log.Error("Failed to send timeout event", "error", err)
		}
		return nil
	case <-ctx.Done():
		log.Info("SSE connection closed", "user_id", userID)
		return nil
	}
}
