package tools

import (
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

func (sp *SchemaProvider) GetResourcesListResource() (string, server.ResourceHandlerFunc) {
	uri := "gorest://resources"
	handler := func() ([]interface{}, error) {
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

		return []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: content,
			},
		}, nil
	}
	return uri, handler
}

func (sp *SchemaProvider) GetSchemaResourceTemplate() (string, server.ResourceTemplateHandlerFunc) {
	uriTemplate := "gorest://schema/{resource}"
	handler := func() (mcp.ResourceTemplate, error) {
		// This will be properly implemented to handle dynamic resource schemas
		return mcp.ResourceTemplate{
			URITemplate: uriTemplate,
			Name:        "Resource Schema",
			Description: "Get schema definition for a specific resource",
			MIMEType:    "application/json",
		}, nil
	}
	return uriTemplate, handler
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
