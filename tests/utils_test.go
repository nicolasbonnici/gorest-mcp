package tests

import (
	"context"
	"testing"

	"github.com/nicolasbonnici/gorest-mcp"
	"github.com/stretchr/testify/assert"
)

// TestGetUserIDFromContext tests user ID extraction from context
func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
		wantID  string
	}{
		{
			name:    "valid user ID",
			ctx:     context.WithValue(context.Background(), mcp.ContextKeyUserID, "user-123"),
			wantErr: false,
			wantID:  "user-123",
		},
		{
			name:    "missing user ID",
			ctx:     context.Background(),
			wantErr: true,
			wantID:  "",
		},
		{
			name:    "empty user ID",
			ctx:     context.WithValue(context.Background(), mcp.ContextKeyUserID, ""),
			wantErr: true,
			wantID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := mcp.GetUserIDFromContext(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, userID)
			}
		})
	}
}

// TestGetUserEmailFromContext tests user email extraction from context
func TestGetUserEmailFromContext(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantErr   bool
		wantEmail string
	}{
		{
			name:      "valid user email",
			ctx:       context.WithValue(context.Background(), mcp.ContextKeyUserEmail, "user@example.com"),
			wantErr:   false,
			wantEmail: "user@example.com",
		},
		{
			name:      "missing user email",
			ctx:       context.Background(),
			wantErr:   true,
			wantEmail: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := mcp.GetUserEmailFromContext(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantEmail, email)
			}
		})
	}
}

// TestContains tests the Contains utility function
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		value    string
		expected bool
	}{
		{
			name:     "value exists",
			slice:    []string{"a", "b", "c"},
			value:    "b",
			expected: true,
		},
		{
			name:     "value does not exist",
			slice:    []string{"a", "b", "c"},
			value:    "d",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			value:    "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mcp.Contains(tt.slice, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToJSON tests JSON marshaling
func TestToJSON(t *testing.T) {
	data := map[string]interface{}{
		"key": "value",
		"num": 42,
	}

	json, err := mcp.ToJSON(data)
	assert.NoError(t, err)
	assert.Contains(t, json, "key")
	assert.Contains(t, json, "value")
}

// TestFromJSON tests JSON unmarshaling
func TestFromJSON(t *testing.T) {
	jsonStr := `{"key":"value","num":42}`

	var result map[string]interface{}
	err := mcp.FromJSON(jsonStr, &result)

	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, float64(42), result["num"])
}
