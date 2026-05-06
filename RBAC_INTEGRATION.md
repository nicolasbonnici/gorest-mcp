# RBAC Integration Guide for gorest-mcp

## Overview

The gorest-mcp plugin **will integrate** with GoREST's RBAC (Role-Based Access Control) system in v0.2.0. This document outlines how RBAC will be enforced for all CRUD operations exposed via MCP tools.

## Current Status (v0.1.0)

- ✅ Authentication: JWT validation implemented
- ✅ User context extraction from tokens
- ⚠️ CRUD operations: Placeholder implementations
- ❌ RBAC enforcement: Not yet implemented (planned for v0.2.0)

## RBAC Architecture

### How GoREST RBAC Works

GoREST enforces RBAC through its **hook system** at multiple levels:

1. **Field-Level Permissions** (via struct tags)
2. **Query-Level Filters** (via ModifySelectQuery hook)
3. **State Validation** (via StateProcessor hook)
4. **Response Filtering** (via Serializer hook)

### Field-Level Permissions

Fields are tagged with `rbac` annotations defining read/write permissions:

```go
type Post struct {
    ID        string    `json:"id" rbac:"read:*"`                        // Everyone can read
    Title     string    `json:"title" rbac:"read:*,write:writer"`        // Writers can edit
    Content   string    `json:"content" rbac:"read:*,write:writer"`      // Writers can edit
    IsDraft   bool      `json:"is_draft" rbac:"read:*,write:writer"`     // Writers can edit
    Secret    string    `json:"secret" rbac:"read:admin,write:admin"`    // Admin only
    CreatedBy string    `json:"created_by" rbac:"read:*,write:none"`     // Read-only
    CreatedAt time.Time `json:"created_at" rbac:"read:*,write:none"`     // Read-only
}
```

**Permission Syntax:**
- `read:*` - Everyone can read
- `read:admin` - Only admin role can read
- `write:writer` - Writer role can write
- `write:none` - No one can write (read-only)
- Multiple roles: `read:writer,admin` or `read:*`

## MCP Tool RBAC Implementation (v0.2.0)

### 1. List Resources (`gorest_list_resources`)

**RBAC Enforcement Points:**

```go
func (ct *CRUDTools) handleListResources(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // 1. Extract user context from JWT (already validated in SSE handler)
    userID := ctx.Value(ContextKeyUserID).(string)
    userRoles, _ := ctx.Value(ContextKeyUserRoles).([]string)

    // 2. Create CRUD instance with hooks
    //    Hooks automatically enforce RBAC
    crudOps := crud.New[models.Post](ct.db, hooks.GetHooks(&models.Post{}))

    // 3. The hook system automatically:
    //    a) ModifySelectQuery: Adds user-scoped filters
    //       - Filters by tenant_id, user_id, or other multi-tenancy fields
    //       - Ensures users only see their own data
    //    b) Serializer: Filters response fields based on read permissions
    //       - Removes fields user doesn't have read access to
    //       - Based on rbac:"read:role" tags

    // 4. Execute query with filters from request
    results, err := crudOps.GetAll(ctx, filters, limit, offset)

    return results, err
}
```

**Example Hook Implementation:**

```go
type PostHooks struct {
    hooks.NoOpHooks[models.Post]
}

// ModifySelectQuery adds user-scoped filters
func (h *PostHooks) ModifySelectQuery(ctx context.Context, op hooks.Operation, builder *query.SelectBuilder) (*query.SelectBuilder, bool) {
    userID := ctx.Value("user_id").(string)
    userRoles := ctx.Value("user_roles").([]string)

    // Non-admin users can only see their own posts
    if !contains(userRoles, "admin") {
        builder = builder.Where(query.Eq("created_by", userID))
    }

    return builder, true
}

// Serializer filters fields based on RBAC
func (h *PostHooks) SerializeMany(ctx context.Context, posts []models.Post) ([]map[string]interface{}, error) {
    userRoles := ctx.Value("user_roles").([]string)

    var results []map[string]interface{}
    for _, post := range posts {
        serialized := make(map[string]interface{})

        // Always include fields with rbac:"read:*"
        serialized["id"] = post.ID
        serialized["title"] = post.Title
        serialized["content"] = post.Content

        // Only include secret for admin
        if contains(userRoles, "admin") {
            serialized["secret"] = post.Secret
        }

        results = append(results, serialized)
    }

    return results, nil
}
```

### 2. Get Resource (`gorest_get_resource`)

**RBAC Enforcement:**

```go
func (ct *CRUDTools) handleGetResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    userID := ctx.Value(ContextKeyUserID).(string)

    // Create CRUD with hooks
    crudOps := crud.New[models.Post](ct.db, hooks.GetHooks(&models.Post{}))

    // GetByID automatically:
    // 1. Checks if user has access (via ModifySelectQuery hook)
    // 2. Returns 404 if not found OR user lacks permission
    // 3. Filters response fields (via Serializer hook)

    result, err := crudOps.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, crud.ErrNotFound) {
            // Could be not found OR permission denied (don't leak info)
            return nil, ErrResourceNotFound
        }
        return nil, err
    }

    return result, nil
}
```

### 3. Create Resource (`gorest_create_resource`)

**RBAC Enforcement:**

```go
func (ct *CRUDTools) handleCreateResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    userID := ctx.Value(ContextKeyUserID).(string)
    userRoles := ctx.Value(ContextKeyUserRoles").([]string)

    // Parse resource data from request
    var resource models.Post
    // ... parse request.Params.Arguments["data"]

    // StateProcessor hook validates:
    // 1. User has write permissions for each field
    // 2. Enriches with user context (created_by, tenant_id)
    // 3. Validates business rules

    crudOps := crud.New[models.Post](ct.db, hooks.GetHooks(&models.Post{}))
    created, err := crudOps.Create(ctx, &resource)

    return created, err
}
```

**Example StateProcessor for Create:**

```go
func (h *PostHooks) StateProcessor(ctx context.Context, op hooks.Operation, id any, post *models.Post) error {
    userID := ctx.Value("user_id").(string)
    userRoles := ctx.Value("user_roles").([]string)

    if op == hooks.OperationCreate {
        // Automatically set created_by
        post.CreatedBy = userID

        // Validate write permissions
        if post.Secret != "" && !contains(userRoles, "admin") {
            return fmt.Errorf("only admins can set secret field")
        }

        if !contains(userRoles, "writer") && !contains(userRoles, "admin") {
            return fmt.Errorf("user lacks writer role")
        }
    }

    return nil
}
```

### 4. Update Resource (`gorest_update_resource`)

**RBAC Enforcement:**

```go
func (ct *CRUDTools) handleUpdateResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    userID := ctx.Value(ContextKeyUserID).(string)
    userRoles := ctx.Value(ContextKeyUserRoles").([]string)

    // Parse updates from request
    var updates map[string]interface{}
    // ... parse request.Params.Arguments["data"]

    // RBAC enforced via:
    // 1. ModifyUpdateQuery: Ensures user can only update their own resources
    // 2. StateProcessor: Validates field write permissions
    // 3. Prevents updating read-only fields (id, created_at, etc.)

    crudOps := crud.New[models.Post](ct.db, hooks.GetHooks(&models.Post{}))
    updated, err := crudOps.Update(ctx, id, updates)

    return updated, err
}
```

**Example ModifyUpdateQuery:**

```go
func (h *PostHooks) ModifyUpdateQuery(ctx context.Context, builder *query.UpdateBuilder) (*query.UpdateBuilder, bool) {
    userID := ctx.Value("user_id").(string)
    userRoles := ctx.Value("user_roles").([]string)

    // Non-admin users can only update their own posts
    if !contains(userRoles, "admin") {
        builder = builder.Where(query.Eq("created_by", userID))
    }

    return builder, true
}
```

### 5. Delete Resource (`gorest_delete_resource`)

**RBAC Enforcement:**

```go
func (ct *CRUDTools) handleDeleteResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    userID := ctx.Value(ContextKeyUserID).(string)

    // ModifyDeleteQuery ensures user can only delete their own resources
    crudOps := crud.New[models.Post](ct.db, hooks.GetHooks(&models.Post{}))
    err := crudOps.Delete(ctx, id)

    if err != nil {
        if errors.Is(err, crud.ErrNotFound) {
            // Could be not found OR permission denied
            return ErrResourceNotFound
        }
        return err
    }

    return nil
}
```

**Soft Delete Example:**

```go
func (h *PostHooks) ModifyDeleteQuery(ctx context.Context, builder *query.DeleteBuilder) (*query.DeleteBuilder, bool) {
    // Return false to skip actual delete, handle in StateProcessor
    return nil, false
}

func (h *PostHooks) StateProcessor(ctx context.Context, op hooks.Operation, id any, post *models.Post) error {
    if op == hooks.OperationDelete {
        userID := ctx.Value("user_id").(string)

        // Check ownership
        if post.CreatedBy != userID {
            return ErrForbidden
        }

        // Soft delete: UPDATE posts SET deleted_at = NOW() WHERE id = ?
        now := time.Now()
        post.DeletedAt = &now

        // Execute UPDATE instead of DELETE
        // ... (would be handled by CRUD layer)
    }

    return nil
}
```

## Schema Introspection with RBAC

The `gorest://schema/{resource}` resource will include RBAC metadata:

```json
{
  "resource": "posts",
  "fields": [
    {
      "name": "id",
      "type": "uuid",
      "rbac_read": ["*"],
      "rbac_write": ["none"]
    },
    {
      "name": "title",
      "type": "string",
      "rbac_read": ["*"],
      "rbac_write": ["writer", "admin"]
    },
    {
      "name": "secret",
      "type": "string",
      "rbac_read": ["admin"],
      "rbac_write": ["admin"]
    }
  ]
}
```

This allows AI agents to:
1. Understand what fields they can access
2. Know which operations are permitted
3. Provide better user experience (hide unauthorized fields)

## Multi-Tenancy Support

RBAC system supports multi-tenancy patterns:

### Pattern 1: User-Scoped Data

```go
// ModifySelectQuery adds user_id filter
builder.Where(query.Eq("user_id", userID))
```

### Pattern 2: Tenant-Scoped Data

```go
// ModifySelectQuery adds tenant_id filter
tenantID := ctx.Value("tenant_id").(string)
builder.Where(query.Eq("tenant_id", tenantID))
```

### Pattern 3: Hierarchical Permissions

```go
// Admin sees all, manager sees team, user sees own
switch {
case contains(userRoles, "admin"):
    // No filters, see everything
case contains(userRoles, "manager"):
    teamID := ctx.Value("team_id").(string)
    builder.Where(query.Eq("team_id", teamID))
default:
    builder.Where(query.Eq("user_id", userID))
}
```

## Security Best Practices

### 1. Never Trust Client Input

Always validate permissions server-side via hooks, never rely on client-provided role claims.

### 2. Fail Securely

When a resource is not found or user lacks permission, return the same error (404) to avoid leaking information about resource existence.

### 3. Audit Logging

Log all CRUD operations with user context for audit trails:

```go
ct.logger.Info("Resource created",
    "user_id", userID,
    "resource", "posts",
    "resource_id", created.ID,
)
```

### 4. Rate Limiting

Apply rate limits per user to prevent abuse (already implemented in middleware).

### 5. Field Sanitization

Always filter response fields via Serializer hook, never expose all fields blindly.

## Testing RBAC Integration

v0.2.0 will include comprehensive tests:

```go
func TestListResourcesRBAC(t *testing.T) {
    // Create posts by different users
    // Login as user1
    // List posts -> should only see user1's posts
    // Login as admin
    // List posts -> should see all posts
}

func TestUpdateResourceRBAC(t *testing.T) {
    // Create post as user1
    // Try to update as user2 -> should fail
    // Try to update as user1 -> should succeed
    // Try to update secret field as user1 -> should fail
    // Try to update secret field as admin -> should succeed
}
```

## Implementation Checklist (v0.2.0)

- [ ] Extract user context (user_id, roles) in all CRUD handlers
- [ ] Implement CRUD operations using GoREST crud package
- [ ] Create hooks for each resource type
- [ ] Implement ModifySelectQuery for read permissions
- [ ] Implement ModifyUpdateQuery for update permissions
- [ ] Implement ModifyDeleteQuery for delete permissions
- [ ] Implement StateProcessor for field-level write validation
- [ ] Implement Serializer for field-level read filtering
- [ ] Add RBAC metadata to schema introspection
- [ ] Write comprehensive RBAC tests
- [ ] Document RBAC patterns and examples
- [ ] Add audit logging for all operations

## Example: Complete RBAC Flow

```
AI Agent Request: "Create a blog post"
  ↓
MCP Tool: gorest_create_resource
  ↓
[1] Extract JWT user context: {user_id: "uuid-123", roles: ["writer"]}
  ↓
[2] Parse resource data: {title: "My Post", content: "...", secret: "admin-only"}
  ↓
[3] StateProcessor Hook:
    ✓ User has "writer" role (can create posts)
    ✓ Auto-set created_by = "uuid-123"
    ✗ User tries to set "secret" field without admin role
    → Return Error: "Only admins can set secret field"
  ↓
[Error returned to AI agent]
```

**Successful Create:**

```
AI Agent Request: "Create a blog post"
  ↓
MCP Tool: gorest_create_resource
  ↓
[1] User context: {user_id: "uuid-123", roles: ["writer"]}
  ↓
[2] Data: {title: "My Post", content: "..."}
  ↓
[3] StateProcessor Hook:
    ✓ User has "writer" role
    ✓ Auto-set created_by = "uuid-123"
    ✓ No admin-only fields set
  ↓
[4] Database INSERT
  ↓
[5] Serializer Hook:
    ✓ Filter response fields based on read permissions
    ✓ User can read: id, title, content, created_by, created_at
    ✗ User cannot read: secret (admin only)
  ↓
[6] Return to AI agent: {id: "uuid-456", title: "My Post", ...}
```

## References

- [GoREST RBAC Documentation](https://github.com/nicolasbonnici/gorest-rbac)
- [GoREST Hooks System](https://github.com/nicolasbonnici/gorest/blob/main/HOOKS.md)
- [Technical Specification](GOREST_MCP_TECHNICAL_SPEC.md) - Section 10.3

---

**Version**: 0.1.0 (Planning Document)
**Target Implementation**: v0.2.0
**Last Updated**: 2026-05-01
