package tools

import (
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest/database"
)

// toolEntry pairs a tool definition with its bound handler for O(1) dispatch.
type toolEntry struct {
	tool    mcp.Tool
	handler server.ToolHandlerFunc
}

type Registry struct {
	db     database.Database
	logger *slog.Logger

	crud   *CRUDTools
	schema *SchemaProvider

	// tools indexes every CRUD tool by name; order preserves a stable
	// registration sequence. Both are built once here so registration and
	// lookups never re-scan or rebuild the tool set on later calls.
	tools map[string]toolEntry
	order []string
}

func NewRegistry(db database.Database, log *slog.Logger) *Registry {
	crud := NewCRUDTools(db, log)
	r := &Registry{
		db:     db,
		logger: log,
		crud:   crud,
		schema: NewSchemaProvider(db, log),
		tools:  make(map[string]toolEntry, 5),
	}

	for _, build := range []func() (mcp.Tool, server.ToolHandlerFunc){
		crud.GetListResourcesTool,
		crud.GetGetResourceTool,
		crud.GetCreateResourceTool,
		crud.GetUpdateResourceTool,
		crud.GetDeleteResourceTool,
	} {
		tool, handler := build()
		r.tools[tool.Name] = toolEntry{tool: tool, handler: handler}
		r.order = append(r.order, tool.Name)
	}

	return r
}

// Tool returns the pre-built tool definition and handler for name in O(1).
func (r *Registry) Tool(name string) (mcp.Tool, server.ToolHandlerFunc, bool) {
	e, ok := r.tools[name]
	return e.tool, e.handler, ok
}

func (r *Registry) RegisterCRUDTools(mcpServer *server.MCPServer) error {
	for _, name := range r.order {
		e := r.tools[name]
		mcpServer.AddTool(e.tool, e.handler)
	}

	r.logger.Info("CRUD tools registered successfully")
	return nil
}

func (r *Registry) RegisterSchemaResources(mcpServer *server.MCPServer) error {
	resource, handler := r.schema.GetResourcesListResource()
	mcpServer.AddResource(resource, handler)

	template, templateHandler := r.schema.GetSchemaResourceTemplate()
	mcpServer.AddResourceTemplate(template, templateHandler)

	r.logger.Info("Schema resources registered successfully")
	return nil
}
