package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest-mcp/tools"
	"github.com/nicolasbonnici/gorest/database"
)

type MCPServer struct {
	server *server.MCPServer
	config *Config
	db     database.Database
	logger *slog.Logger
	tools  *tools.Registry
}

func NewMCPServer(config *Config, db database.Database, log *slog.Logger) (*MCPServer, error) {
	mcpServer := server.NewMCPServer(
		"GoREST MCP Server",
		"0.1.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, false),
	)

	// Initialize tool registry
	registry := tools.NewRegistry(db, log)

	// Create wrapper
	wrapper := &MCPServer{
		server: mcpServer,
		config: config,
		db:     db,
		logger: log,
		tools:  registry,
	}

	// Register tools and resources
	if err := wrapper.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	if err := wrapper.registerResources(); err != nil {
		return nil, fmt.Errorf("failed to register resources: %w", err)
	}

	log.Info("MCP server initialized with tools and resources")
	return wrapper, nil
}

// registerTools registers MCP tools based on enabled operations
func (s *MCPServer) registerTools() error {
	// Register CRUD tools if enabled
	if s.config.IsOperationEnabled("crud") {
		if err := s.tools.RegisterCRUDTools(s.server); err != nil {
			return fmt.Errorf("failed to register CRUD tools: %w", err)
		}
		s.logger.Info("CRUD tools registered")
	}

	return nil
}

// registerResources registers MCP resources based on enabled operations
func (s *MCPServer) registerResources() error {
	// Register schema resources if enabled
	if s.config.IsOperationEnabled("schema") {
		if err := s.tools.RegisterSchemaResources(s.server); err != nil {
			return fmt.Errorf("failed to register schema resources: %w", err)
		}
		s.logger.Info("Schema resources registered")
	}

	return nil
}

// HandleRequest processes an MCP request with user context
func (s *MCPServer) HandleRequest(ctx context.Context, request interface{}) (interface{}, error) {
	// Extract user context for logging
	userID, _ := GetUserIDFromContext(ctx)
	if userID != "" && s.config.LogRequests {
		s.logger.Info("Processing MCP request", "user_id", userID)
	}

	// Forward to underlying MCP server
	// Note: The actual request handling will be done via SSE transport
	return nil, nil
}

// GetServer returns the underlying MCP server for SSE transport
func (s *MCPServer) GetServer() *server.MCPServer {
	return s.server
}
