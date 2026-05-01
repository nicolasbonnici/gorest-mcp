package tests

import (
	"errors"
	"testing"

	"github.com/nicolasbonnici/gorest-mcp"
	"github.com/stretchr/testify/assert"
)

// TestWrapError tests error wrapping functionality
func TestWrapError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "unauthorized error",
			err:          mcp.ErrUnauthorized,
			expectedCode: mcp.CodeUnauthorized,
		},
		{
			name:         "forbidden error",
			err:          mcp.ErrForbidden,
			expectedCode: mcp.CodeForbidden,
		},
		{
			name:         "rate limit exceeded",
			err:          mcp.ErrRateLimitExceeded,
			expectedCode: mcp.CodeRateLimitExceeded,
		},
		{
			name:         "resource not found",
			err:          mcp.ErrResourceNotFound,
			expectedCode: mcp.CodeResourceNotFound,
		},
		{
			name:         "operation disabled",
			err:          mcp.ErrOperationDisabled,
			expectedCode: mcp.CodeOperationDisabled,
		},
		{
			name:         "max connections exceeded",
			err:          mcp.ErrMaxConnectionsExceeded,
			expectedCode: mcp.CodeMaxConnectionsExceeded,
		},
		{
			name:         "connection timeout",
			err:          mcp.ErrConnectionTimeout,
			expectedCode: mcp.CodeConnectionTimeout,
		},
		{
			name:         "invalid request",
			err:          mcp.ErrInvalidRequest,
			expectedCode: mcp.CodeInvalidRequest,
		},
		{
			name:         "generic error",
			err:          errors.New("some generic error"),
			expectedCode: mcp.CodeInternalError,
		},
		{
			name:         "nil error",
			err:          nil,
			expectedCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcpErr := mcp.WrapError(tt.err)
			if tt.err == nil {
				assert.Nil(t, mcpErr)
			} else {
				assert.NotNil(t, mcpErr)
				assert.Equal(t, tt.expectedCode, mcpErr.Code)
			}
		})
	}
}

// TestMCPErrorUnwrap tests error unwrapping
func TestMCPErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	mcpErr := mcp.NewMCPError(mcp.CodeInternalError, "wrapped error", originalErr)

	assert.True(t, errors.Is(mcpErr, originalErr))
}

// TestMCPErrorAlreadyWrapped tests that wrapping an MCPError returns it as-is
func TestMCPErrorAlreadyWrapped(t *testing.T) {
	originalMCPErr := mcp.NewMCPError(mcp.CodeUnauthorized, "unauthorized", nil)
	wrappedErr := mcp.WrapError(originalMCPErr)

	assert.Equal(t, originalMCPErr, wrappedErr)
}
