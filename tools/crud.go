package tools

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nicolasbonnici/gorest/database"
)

type CRUDTools struct {
	db     database.Database
	logger *slog.Logger
}

func NewCRUDTools(db database.Database, log *slog.Logger) *CRUDTools {
	return &CRUDTools{
		db:     db,
		logger: log,
	}
}

func (ct *CRUDTools) GetListResourcesTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "gorest_list_resources",
		Description: "List all resources with pagination and filtering",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"resource": map[string]interface{}{
					"type":        "string",
					"description": "Resource name (e.g., 'posts', 'users')",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Items per page (default: 20)",
					"default":     20,
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Offset for pagination (default: 0)",
					"default":     0,
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": "Filter conditions (key-value pairs)",
				},
				"order_by": map[string]interface{}{
					"type":        "string",
					"description": "Sort field",
				},
				"order": map[string]interface{}{
					"type":        "string",
					"description": "Sort order: asc or desc (default: asc)",
					"enum":        []string{"asc", "desc"},
					"default":     "asc",
				},
			},
		},
	}
	return tool, ct.handleListResources
}

func (ct *CRUDTools) GetGetResourceTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "gorest_get_resource",
		Description: "Get a single resource by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"resource": map[string]interface{}{
					"type":        "string",
					"description": "Resource name (e.g., 'posts', 'users')",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Resource ID",
				},
			},
		},
	}
	return tool, ct.handleGetResource
}

func (ct *CRUDTools) GetCreateResourceTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "gorest_create_resource",
		Description: "Create a new resource",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"resource": map[string]interface{}{
					"type":        "string",
					"description": "Resource name (e.g., 'posts', 'users')",
				},
				"data": map[string]interface{}{
					"type":        "object",
					"description": "Resource data to create",
				},
			},
		},
	}
	return tool, ct.handleCreateResource
}

func (ct *CRUDTools) GetUpdateResourceTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "gorest_update_resource",
		Description: "Update an existing resource",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"resource": map[string]interface{}{
					"type":        "string",
					"description": "Resource name (e.g., 'posts', 'users')",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Resource ID",
				},
				"data": map[string]interface{}{
					"type":        "object",
					"description": "Resource data to update",
				},
			},
		},
	}
	return tool, ct.handleUpdateResource
}

func (ct *CRUDTools) GetDeleteResourceTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.Tool{
		Name:        "gorest_delete_resource",
		Description: "Delete a resource by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"resource": map[string]interface{}{
					"type":        "string",
					"description": "Resource name (e.g., 'posts', 'users')",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Resource ID",
				},
			},
		},
	}
	return tool, ct.handleDeleteResource
}

// Tool handlers - TODO: implement full CRUD with RBAC in separate task

func (ct *CRUDTools) handleListResources(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	ct.logger.Info("handleListResources called", "arguments", args)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("List resources tool called with arguments: %v\n\nNote: Full CRUD implementation with RBAC coming soon", args),
			},
		},
	}, nil
}

func (ct *CRUDTools) handleGetResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	ct.logger.Info("handleGetResource called", "arguments", args)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Get resource tool called with arguments: %v\n\nNote: Full CRUD implementation with RBAC coming soon", args),
			},
		},
	}, nil
}

func (ct *CRUDTools) handleCreateResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	ct.logger.Info("handleCreateResource called", "arguments", args)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Create resource tool called with arguments: %v\n\nNote: Full CRUD implementation with RBAC coming soon", args),
			},
		},
	}, nil
}

func (ct *CRUDTools) handleUpdateResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	ct.logger.Info("handleUpdateResource called", "arguments", args)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Update resource tool called with arguments: %v\n\nNote: Full CRUD implementation with RBAC coming soon", args),
			},
		},
	}, nil
}

func (ct *CRUDTools) handleDeleteResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	ct.logger.Info("handleDeleteResource called", "arguments", args)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Delete resource tool called with arguments: %v\n\nNote: Full CRUD implementation with RBAC coming soon", args),
			},
		},
	}, nil
}
