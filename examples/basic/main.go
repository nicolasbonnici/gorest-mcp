package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolasbonnici/gorest"
	"github.com/nicolasbonnici/gorest/database"
	mcp "github.com/nicolasbonnici/gorest-mcp"
)

func main() {
	// Initialize GoREST application
	app, err := gorest.New(&gorest.Config{
		AppName: "GoREST MCP Example",
		Port:    8080,
		Database: database.Config{
			Driver: "sqlite",
			DSN:    "gorest_mcp_example.db",
		},
		Plugins: []gorest.PluginConfig{
			{
				Name:    "mcp",
				Enabled: true,
				Config: map[string]interface{}{
					"enabled": true,
					"enabled_operations": []string{"crud", "schema"},
					"rate_limit": map[string]interface{}{
						"enabled":             true,
						"requests_per_minute": 60,
						"burst":               10,
					},
					"sse": map[string]interface{}{
						"heartbeat_interval":      30,
						"connection_timeout":      300,
						"max_connections_per_user": 3,
					},
					"log_requests": true,
					"log_level":    "info",
				},
			},
		},
	})

	if err != nil {
		log.Fatal("Failed to initialize GoREST:", err)
	}

	// Register MCP plugin
	mcpPlugin := &mcp.Plugin{}
	if err := app.RegisterPlugin(mcpPlugin); err != nil {
		log.Fatal("Failed to register MCP plugin:", err)
	}

	// Add health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "gorest-mcp-example",
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting GoREST MCP Example on port %s", port)
	log.Printf("MCP endpoint available at: http://localhost:%s/mcp", port)
	log.Fatal(app.Listen(":" + port))
}
