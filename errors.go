package mcp

import (
	"errors"
	"fmt"
)

var (
	// ErrUnauthorized indicates missing or invalid authentication
	ErrUnauthorized = errors.New("unauthorized: invalid or missing JWT token")

	// ErrForbidden indicates insufficient permissions
	ErrForbidden = errors.New("forbidden: insufficient permissions")

	// ErrRateLimitExceeded indicates rate limit has been exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidRequest indicates malformed request
	ErrInvalidRequest = errors.New("invalid request")

	// ErrResourceNotFound indicates requested resource doesn't exist
	ErrResourceNotFound = errors.New("resource not found")

	// ErrOperationDisabled indicates the requested operation is disabled
	ErrOperationDisabled = errors.New("operation disabled")

	// ErrMaxConnectionsExceeded indicates user has too many active connections
	ErrMaxConnectionsExceeded = errors.New("maximum connections exceeded")

	// ErrConnectionTimeout indicates connection idle timeout
	ErrConnectionTimeout = errors.New("connection timeout")
)

// MCPError wraps an error with additional context for MCP protocol
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *MCPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work
func (e *MCPError) Unwrap() error {
	return e.Err
}

// NewMCPError creates a new MCP error
func NewMCPError(code int, message string, err error) *MCPError {
	return &MCPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error code constants following JSON-RPC 2.0 spec
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603

	// Custom error codes (application-specific)
	CodeUnauthorized           = -32000
	CodeForbidden              = -32001
	CodeRateLimitExceeded      = -32002
	CodeResourceNotFound       = -32003
	CodeOperationDisabled      = -32004
	CodeMaxConnectionsExceeded = -32005
	CodeConnectionTimeout      = -32006
)

// WrapError wraps a standard error into an MCPError
func WrapError(err error) *MCPError {
	if err == nil {
		return nil
	}

	// If already an MCPError, return as-is
	var mcpErr *MCPError
	if errors.As(err, &mcpErr) {
		return mcpErr
	}

	// Map known errors to MCP errors
	switch {
	case errors.Is(err, ErrUnauthorized):
		return NewMCPError(CodeUnauthorized, "Unauthorized", err)
	case errors.Is(err, ErrForbidden):
		return NewMCPError(CodeForbidden, "Forbidden", err)
	case errors.Is(err, ErrRateLimitExceeded):
		return NewMCPError(CodeRateLimitExceeded, "Rate limit exceeded", err)
	case errors.Is(err, ErrResourceNotFound):
		return NewMCPError(CodeResourceNotFound, "Resource not found", err)
	case errors.Is(err, ErrOperationDisabled):
		return NewMCPError(CodeOperationDisabled, "Operation disabled", err)
	case errors.Is(err, ErrMaxConnectionsExceeded):
		return NewMCPError(CodeMaxConnectionsExceeded, "Maximum connections exceeded", err)
	case errors.Is(err, ErrConnectionTimeout):
		return NewMCPError(CodeConnectionTimeout, "Connection timeout", err)
	case errors.Is(err, ErrInvalidRequest):
		return NewMCPError(CodeInvalidRequest, "Invalid request", err)
	default:
		return NewMCPError(CodeInternalError, "Internal error", err)
	}
}
