package tests

import (
	"testing"

	"github.com/nicolasbonnici/gorest-mcp"
	"github.com/stretchr/testify/assert"
)

// TestPluginName verifies the plugin name is correct
func TestPluginName(t *testing.T) {
	plugin := &mcp.Plugin{}
	assert.Equal(t, "mcp", plugin.Name())
}

// TestDefaultConfig verifies default configuration values
func TestDefaultConfig(t *testing.T) {
	config := mcp.DefaultConfig()

	assert.True(t, config.Enabled)
	assert.Contains(t, config.EnabledOperations, "crud")
	assert.Contains(t, config.EnabledOperations, "schema")
	assert.True(t, config.RateLimit.Enabled)
	assert.Equal(t, 60, config.RateLimit.RequestsPerMinute)
	assert.Equal(t, 10, config.RateLimit.Burst)
	assert.Equal(t, 30, config.SSE.HeartbeatInterval)
	assert.Equal(t, 300, config.SSE.ConnectionTimeout)
	assert.Equal(t, 3, config.SSE.MaxConnectionsPerUser)
	assert.True(t, config.LogRequests)
	assert.Equal(t, "info", config.LogLevel)
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *mcp.Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  mcp.DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid operation",
			config: &mcp.Config{
				Enabled:           true,
				EnabledOperations: []string{"invalid"},
				RateLimit: mcp.RateLimitConfig{
					Enabled:           true,
					RequestsPerMinute: 60,
					Burst:             10,
				},
				SSE: mcp.SSEConfig{
					HeartbeatInterval:     30,
					ConnectionTimeout:     300,
					MaxConnectionsPerUser: 3,
				},
				LogLevel: "info",
			},
			wantErr: true,
		},
		{
			name: "invalid rate limit",
			config: &mcp.Config{
				Enabled:           true,
				EnabledOperations: []string{"crud"},
				RateLimit: mcp.RateLimitConfig{
					Enabled:           true,
					RequestsPerMinute: 0, // Invalid
					Burst:             10,
				},
				SSE: mcp.SSEConfig{
					HeartbeatInterval:     30,
					ConnectionTimeout:     300,
					MaxConnectionsPerUser: 3,
				},
				LogLevel: "info",
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &mcp.Config{
				Enabled:           true,
				EnabledOperations: []string{"crud"},
				RateLimit: mcp.RateLimitConfig{
					Enabled:           true,
					RequestsPerMinute: 60,
					Burst:             10,
				},
				SSE: mcp.SSEConfig{
					HeartbeatInterval:     30,
					ConnectionTimeout:     300,
					MaxConnectionsPerUser: 3,
				},
				LogLevel: "invalid", // Invalid
			},
			wantErr: true,
		},
		{
			name: "disabled plugin",
			config: &mcp.Config{
				Enabled: false, // Should pass validation even with invalid settings
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestConfigOperationEnabled tests the IsOperationEnabled method
func TestConfigOperationEnabled(t *testing.T) {
	config := mcp.DefaultConfig()

	assert.True(t, config.IsOperationEnabled("crud"))
	assert.True(t, config.IsOperationEnabled("schema"))
	assert.False(t, config.IsOperationEnabled("nonexistent"))
}

// TestParseConfig tests configuration parsing
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "nil config returns default",
			input:   nil,
			wantErr: false,
		},
		{
			name: "valid map config",
			input: map[string]interface{}{
				"enabled": true,
				"enabled_operations": []interface{}{"crud", "schema"},
				"log_level": "debug",
			},
			wantErr: false,
		},
		{
			name:    "valid Config struct",
			input:   mcp.DefaultConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := mcp.ParseConfig(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			}
		})
	}
}
