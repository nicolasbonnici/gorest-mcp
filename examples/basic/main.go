package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nicolasbonnici/gorest"
	"github.com/nicolasbonnici/gorest/database"
	"github.com/nicolasbonnici/gorest/plugin"
	"github.com/nicolasbonnici/gorest/pluginloader"

	mcpplugin "github.com/nicolasbonnici/gorest-mcp"
)

func init() {
	// Register MCP plugin factory
	pluginloader.RegisterPluginFactory("mcp", mcpplugin.NewPlugin)
}

func registerRoutes(router fiber.Router, db database.Database, paginationLimit, paginationMaxLimit int, pluginRegistry *plugin.PluginRegistry) {
	// Custom routes can be registered here
	router.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "gorest-mcp-example",
		})
	})
}

func main() {
	cfg := gorest.Config{
		ConfigPath:     ".",
		RegisterRoutes: registerRoutes,
	}

	gorest.Start(cfg)
}
