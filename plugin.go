package mcp

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolasbonnici/gorest/database"
	"github.com/nicolasbonnici/gorest/logger"
	"github.com/nicolasbonnici/gorest/plugin"
)

type Plugin struct {
	config    *Config
	db        database.Database
	logger    *slog.Logger
	mcpServer *MCPServer
}

func (p *Plugin) Name() string {
	return "mcp"
}

func (p *Plugin) Initialize(cfg map[string]any) error {
	// Extract database from shared config
	if db, ok := cfg["database"].(database.Database); ok {
		p.db = db
	} else {
		return fmt.Errorf("database not found in config")
	}

	// Use GoREST global logger
	p.logger = logger.Log

	// Extract plugin-specific config
	var pluginConfig interface{}
	if config, ok := cfg["config"]; ok {
		pluginConfig = config
	}

	// Parse and validate configuration
	config, err := ParseConfig(pluginConfig)
	if err != nil {
		return fmt.Errorf("failed to parse MCP config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid MCP config: %w", err)
	}

	p.config = config

	if !config.Enabled {
		p.logger.Info("MCP plugin is disabled")
		return nil
	}

	// Initialize MCP server
	mcpServer, err := NewMCPServer(config, p.db, p.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize MCP server: %w", err)
	}

	p.mcpServer = mcpServer

	p.logger.Info("MCP plugin initialized successfully")
	return nil
}

func (p *Plugin) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func (p *Plugin) SetupEndpoints(router fiber.Router) error {
	if p.config == nil || !p.config.Enabled {
		return nil
	}

	router.Get("/mcp", p.handleSSE)

	p.logger.Info("MCP endpoints registered: GET /mcp")
	return nil
}

func (p *Plugin) handleSSE(c *fiber.Ctx) error {
	return HandleSSE(c, p.mcpServer, p.config, p.logger)
}

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.EndpointSetup = (*Plugin)(nil)
