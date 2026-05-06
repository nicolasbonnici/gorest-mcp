package tools

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest/database"
)

type SchemaProvider struct {
	db     database.Database
	logger *slog.Logger
}

func NewSchemaProvider(db database.Database, log *slog.Logger) *SchemaProvider {
	return &SchemaProvider{
		db:     db,
		logger: log,
	}
}

func (sp *SchemaProvider) GetResourcesListResource() (mcp.Resource, server.ResourceHandlerFunc) {
	resource := mcp.NewResource(
		"gorest://resources",
		"Resources List",
		func(r *mcp.Resource) {
			r.Description = "List of all available GoREST resources"
			r.MIMEType = "application/json"
		},
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		sp.logger.Info("ResourcesList called")

		content := `{
  "resources": [
    {
      "name": "example",
      "table": "examples",
      "description": "Example resource for demonstration",
      "endpoints": [
        "GET /api/examples",
        "POST /api/examples",
        "GET /api/examples/:id",
        "PUT /api/examples/:id",
        "DELETE /api/examples/:id"
      ]
    }
  ],
  "note": "Full schema introspection coming soon"
}`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "gorest://resources",
				MIMEType: "application/json",
				Text:     content,
			},
		}, nil
	}
	return resource, handler
}

func (sp *SchemaProvider) GetSchemaResourceTemplate() (mcp.ResourceTemplate, server.ResourceTemplateHandlerFunc) {
	template := mcp.NewResourceTemplate(
		"gorest://schema/{resource}",
		"Resource Schema",
		func(t *mcp.ResourceTemplate) {
			t.Description = "Get schema definition for a specific resource"
			t.MIMEType = "application/json"
		},
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		sp.logger.Info("GetSchema called", "uri", request.Params.URI)

		// Extract resource name from URI
		// For now, return a placeholder - full implementation in v0.2.0
		content := `{
  "resource": "example",
  "table": "examples",
  "fields": [
    {
      "name": "id",
      "type": "uuid",
      "nullable": false,
      "primary_key": true,
      "auto_generated": true
    },
    {
      "name": "name",
      "type": "string",
      "max_length": 255,
      "nullable": false,
      "validation": "required,min=1,max=255"
    },
    {
      "name": "created_at",
      "type": "timestamp",
      "nullable": false,
      "auto_generated": true
    }
  ],
  "note": "Full schema introspection coming soon"
}`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     content,
			},
		}, nil
	}
	return template, handler
}

func (sp *SchemaProvider) GetSchemaForResource(resource string) ([]interface{}, error) {
	sp.logger.Info("GetSchema called", "resource", resource)

	content := fmt.Sprintf(`{
  "resource": "%s",
  "table": "%s",
  "fields": [
    {
      "name": "id",
      "type": "uuid",
      "nullable": false,
      "primary_key": true,
      "auto_generated": true
    },
    {
      "name": "name",
      "type": "string",
      "max_length": 255,
      "nullable": false,
      "validation": "required,min=1,max=255"
    },
    {
      "name": "created_at",
      "type": "timestamp",
      "nullable": false,
      "auto_generated": true
    }
  ],
  "note": "Full schema introspection coming soon"
}`, resource, resource)

	return []interface{}{
		mcp.TextContent{
			Type: "text",
			Text: content,
		},
	}, nil
}
