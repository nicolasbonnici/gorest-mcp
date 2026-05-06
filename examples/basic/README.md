# GoREST MCP - Basic Example

This example demonstrates how to integrate the MCP plugin into a GoREST application.

## Setup

1. **Install dependencies:**
```bash
go mod download
```

2. **Create environment file:**
```bash
cp .env.example .env
# Edit .env and set your JWT_SECRET (minimum 32 characters)
```

3. **Run the example:**
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in `gorest.yaml` or `PORT` environment variable).

## Configuration

The example uses `gorest.yaml` for configuration. Key MCP settings:

- **enabled_operations**: List of enabled operations (`crud`, `schema`)
- **rate_limit**: Rate limiting configuration
- **sse**: Server-Sent Events configuration
- **log_level**: Logging level (`debug`, `info`, `warn`, `error`)

See `gorest.yaml` for full configuration options.

## Usage

### 1. Authenticate

First, create a user and get a JWT token:

```bash
# Register a new user (if registration is enabled)
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Login to get JWT token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response: {"token":"eyJhbGc...","user":{...}}
```

### 2. Connect to MCP

The MCP server is available at the `/mcp` endpoint via Server-Sent Events (SSE).

**Using Claude Desktop:**

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "gorest-example": {
      "url": "http://localhost:8080/mcp",
      "transport": "sse",
      "headers": {
        "Authorization": "Bearer YOUR_JWT_TOKEN_HERE"
      }
    }
  }
}
```

Replace `YOUR_JWT_TOKEN_HERE` with the token from the login response.

### 3. Use MCP Tools

Once connected, AI agents can use these tools:

- **gorest_list_resources** - List resources with filtering/pagination
- **gorest_get_resource** - Get a single resource by ID
- **gorest_create_resource** - Create a new resource
- **gorest_update_resource** - Update an existing resource
- **gorest_delete_resource** - Delete a resource

And access these resources:

- **gorest://resources** - List all available API resources
- **gorest://schema/{resource}** - Get schema for a specific resource

### 4. Health Check

Check if the server is running:

```bash
curl http://localhost:8080/health
```

## Troubleshooting

### "JWT_SECRET must be at least 32 characters"

Make sure your `.env` file contains a JWT_SECRET with at least 32 characters:

```bash
JWT_SECRET=your-secret-key-here-must-be-at-least-32-characters-long
```

### "Failed to connect to database"

Check that the database URL is correct in `gorest.yaml` or `DATABASE_URL` environment variable.

For SQLite (default):
```yaml
database:
  url: "sqlite://gorest_mcp_example.db"
```

### "Unauthorized"

Make sure you're including the JWT token in the Authorization header:

```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

## Development

### Enable Debug Logging

Update `gorest.yaml`:

```yaml
plugins:
- name: mcp
  config:
    log_level: "debug"
    log_requests: true
```

### Customize Rate Limits

Update `gorest.yaml`:

```yaml
plugins:
- name: mcp
  config:
    rate_limit:
      enabled: true
      requests_per_minute: 100  # Increase limit
      burst: 20                 # Increase burst
```

## Next Steps

- Explore the [full documentation](../../README.md)
- Check out the [technical specification](../../GOREST_MCP_TECHNICAL_SPEC.md)
- Review [CLAUDE.md](../../CLAUDE.md) for AI assistant guidance
