# GoREST-MCP Plugin

[![Go Version](https://img.shields.io/badge/go-1.23-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0-orange.svg)](https://github.com/nicolasbonnici/gorest-mcp/releases)

The **gorest-mcp** plugin integrates the [Model Context Protocol (MCP)](https://modelcontextprotocol.io) into GoREST applications, enabling AI agents (like Claude) to interact with GoREST APIs through a standardized protocol.

## Features

- **MCP Server Embedded**: Full MCP server as a GoREST plugin endpoint
- **SSE Transport**: Server-Sent Events for real-time AI agent communication
- **JWT Authentication**: Secure access using existing GoREST authentication
- **CRUD Operations**: Expose database operations as MCP tools
- **Schema Introspection**: Runtime API discovery for AI agents
- **Multi-Database Support**: PostgreSQL, MySQL, SQLite

## Installation

```bash
go get github.com/nicolasbonnici/gorest-mcp@v0.1.0
```

## Quick Start

### 1. Add MCP Plugin to Your GoREST Application

```go
package main

import (
    "log"

    "github.com/nicolasbonnici/gorest"
    mcp "github.com/nicolasbonnici/gorest-mcp"
)

func main() {
    app, err := gorest.New(&gorest.Config{
        AppName: "My API",
        Port:    8080,
        Plugins: []gorest.PluginConfig{
            {
                Name:    "mcp",
                Enabled: true,
                Config: map[string]interface{}{
                    "enabled": true,
                    "enabled_operations": []string{"crud", "schema"},
                },
            },
        },
    })

    if err != nil {
        log.Fatal(err)
    }

    // Register MCP plugin
    mcpPlugin := &mcp.Plugin{}
    app.RegisterPlugin(mcpPlugin)

    // Start server
    app.Listen(":8080")
}
```

### 2. Authenticate and Get JWT Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Response: {"token":"eyJhbGc...","user":{...}}
```

### 3. Connect AI Agent via SSE

Configure your MCP client (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "my-api": {
      "url": "http://localhost:8080/mcp",
      "transport": "sse",
      "headers": {
        "Authorization": "Bearer YOUR_JWT_TOKEN"
      }
    }
  }
}
```

### 4. Use MCP Tools

Once connected, AI agents can use these tools:

- **gorest_list_resources** - List resources with filtering/pagination
- **gorest_get_resource** - Get a single resource by ID
- **gorest_create_resource** - Create a new resource
- **gorest_update_resource** - Update an existing resource
- **gorest_delete_resource** - Delete a resource

And access these resources:

- **gorest://resources** - List all available API resources
- **gorest://schema/{resource}** - Get schema for a specific resource

## Configuration

### Full Configuration Example

```yaml
plugins:
  - name: mcp
    enabled: true
    config:
      # Features
      enabled_operations:
        - crud                         # CRUD tools
        - schema                       # Schema introspection

      # Rate limiting (requests per minute)
      rate_limit:
        enabled: true
        requests_per_minute: 60
        burst: 10

      # SSE settings
      sse:
        heartbeat_interval: 30         # Seconds between heartbeat events
        connection_timeout: 300        # Seconds before idle connection closes
        max_connections_per_user: 3    # Max concurrent SSE connections

      # Logging
      log_requests: true
      log_level: "info"                # debug, info, warn, error
```

### Environment Variables

```bash
# Required (inherited from GoREST auth)
JWT_SECRET=your-secret-key-min-32-chars

# Optional (plugin-specific)
MCP_ENABLED=true
MCP_RATE_LIMIT=60
MCP_LOG_LEVEL=info
```

## Architecture

```
┌─────────────────────────────────────────┐
│       GoREST Application                │
├─────────────────────────────────────────┤
│  Core Middleware                         │
│  ├── Security, CORS, Logger              │
│  └── Content Negotiation                 │
├─────────────────────────────────────────┤
│  Plugin Middleware                       │
│  ├── Auth Plugin (JWT validation)        │
│  └── MCP Plugin                          │
├─────────────────────────────────────────┤
│  Routes                                  │
│  ├── /api/v1/resources (standard CRUD)   │
│  ├── /auth/* (login, register)           │
│  └── /mcp (SSE endpoint) ─────┐          │
└────────────────────────────────┼──────────┘
                                 │
                                 ▼
                    ┌────────────────────┐
                    │   MCP Protocol     │
                    │   ├── SSE Transport│
                    │   ├── Tool Registry│
                    │   └── Resources    │
                    └────────────────────┘
                                 │
                                 ▼
                        ┌────────────┐
                        │ AI Agent   │
                        │ (Claude)   │
                        └────────────┘
```

## MCP Tools Reference

### gorest_list_resources

List all resources with pagination and filtering.

**Parameters:**
```json
{
  "resource": "posts",       // Required: resource name
  "limit": 20,               // Optional: items per page (default: 20)
  "offset": 0,               // Optional: offset (default: 0)
  "filters": {               // Optional: filter conditions
    "status": "published"
  },
  "order_by": "created_at",  // Optional: sort field
  "order": "desc"            // Optional: asc/desc
}
```

### gorest_get_resource

Get a single resource by ID.

**Parameters:**
```json
{
  "resource": "posts",       // Required: resource name
  "id": "uuid-or-id"         // Required: resource ID
}
```

### gorest_create_resource

Create a new resource.

**Parameters:**
```json
{
  "resource": "posts",       // Required: resource name
  "data": {                  // Required: resource data
    "title": "New Post",
    "content": "..."
  }
}
```

### gorest_update_resource

Update an existing resource.

**Parameters:**
```json
{
  "resource": "posts",       // Required: resource name
  "id": "uuid-here",         // Required: resource ID
  "data": {                  // Required: update data
    "title": "Updated Title"
  }
}
```

### gorest_delete_resource

Delete a resource.

**Parameters:**
```json
{
  "resource": "posts",       // Required: resource name
  "id": "uuid-here"          // Required: resource ID
}
```

## MCP Resources Reference

### gorest://resources

List all available resources in the API.

**Response:**
```json
{
  "resources": [
    {
      "name": "posts",
      "table": "posts",
      "description": "Blog posts",
      "endpoints": ["GET /api/posts", "POST /api/posts", ...]
    }
  ]
}
```

### gorest://schema/{resource}

Get schema definition for a resource.

**Response:**
```json
{
  "resource": "posts",
  "table": "posts",
  "fields": [
    {
      "name": "id",
      "type": "uuid",
      "nullable": false,
      "primary_key": true
    }
  ]
}
```

## Security

### Authentication Flow

1. **Login**: Get JWT token from `/auth/login`
2. **Connect**: Pass token via `Authorization: Bearer <token>` header
3. **Validate**: MCP plugin validates JWT on every request
4. **Authorize**: RBAC enforced via GoREST hooks

### Best Practices

- **Use HTTPS** in production for SSE connections
- **Rotate JWT secrets** regularly (minimum 32 characters)
- **Set short token expiry** (default: 1 hour, configurable)
- **Enable rate limiting** to prevent abuse
- **Monitor connections** using logs and metrics

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run Example

```bash
cd examples/basic
go run main.go
```

### Run with Docker

```bash
docker build -t gorest-mcp .
docker run -p 8080:8080 -e JWT_SECRET=your-secret gorest-mcp
```

## Version 0.1.0 Notes

This is the initial release with core functionality:

- ✅ Plugin interface implementation
- ✅ SSE transport with JWT authentication
- ✅ MCP server initialization with mark3labs/mcp-go
- ✅ Tool registration (CRUD placeholders)
- ✅ Resource registration (schema placeholders)
- ✅ Configuration system
- ✅ Rate limiting
- ✅ Connection management

**Coming in v0.2.0:**

- Full CRUD operations with database integration
- Complete schema introspection
- Enhanced error handling
- Comprehensive tests
- Performance optimizations

## Examples

See the [`examples/`](./examples/) directory for:

- **basic/** - Basic GoREST app with MCP plugin
- **claude-desktop/** - Claude Desktop configuration

## Documentation

- **[Technical Specification](GOREST_MCP_TECHNICAL_SPEC.md)** - Complete technical details
- **[CLAUDE.md](CLAUDE.md)** - AI assistant guidance
- **[GoREST Documentation](https://github.com/nicolasbonnici/gorest)** - Core framework docs

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details

## Links

- **GoREST Framework**: https://github.com/nicolasbonnici/gorest
- **MCP Specification**: https://modelcontextprotocol.io
- **mark3labs/mcp-go**: https://github.com/mark3labs/mcp-go

## Support

- **Issues**: https://github.com/nicolasbonnici/gorest-mcp/issues
- **Discussions**: https://github.com/nicolasbonnici/gorest-mcp/discussions

## Acknowledgments

Built with:
- [GoREST Framework](https://github.com/nicolasbonnici/gorest)
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- [Fiber](https://github.com/gofiber/fiber)
