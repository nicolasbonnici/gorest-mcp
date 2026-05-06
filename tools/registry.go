package tools

import (
	"log/slog"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest/database"
)

type Registry struct {
	db     database.Database
	logger *slog.Logger
}

func NewRegistry(db database.Database, log *slog.Logger) *Registry {
	return &Registry{
		db:     db,
		logger: log,
	}
}

func (r *Registry) RegisterCRUDTools(mcpServer *server.MCPServer) error {
	crudTools := NewCRUDTools(r.db, r.logger)

	tool, handler := crudTools.GetListResourcesTool()
	mcpServer.AddTool(tool, handler)

	tool, handler = crudTools.GetGetResourceTool()
	mcpServer.AddTool(tool, handler)

	tool, handler = crudTools.GetCreateResourceTool()
	mcpServer.AddTool(tool, handler)

	tool, handler = crudTools.GetUpdateResourceTool()
	mcpServer.AddTool(tool, handler)

	tool, handler = crudTools.GetDeleteResourceTool()
	mcpServer.AddTool(tool, handler)

	r.logger.Info("CRUD tools registered successfully")
	return nil
}

func (r *Registry) RegisterSchemaResources(mcpServer *server.MCPServer) error {
	schemaProvider := NewSchemaProvider(r.db, r.logger)

	resource, handler := schemaProvider.GetResourcesListResource()
	mcpServer.AddResource(resource, handler)

	template, templateHandler := schemaProvider.GetSchemaResourceTemplate()
	mcpServer.AddResourceTemplate(template, templateHandler)

	r.logger.Info("Schema resources registered successfully")
	return nil
}
