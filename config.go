package mcp

import (
	"fmt"
	"time"
)

// Config holds the MCP plugin configuration
type Config struct {
	Enabled           bool              `json:"enabled" yaml:"enabled"`
	EnabledOperations []string          `json:"enabled_operations" yaml:"enabled_operations"`
	RateLimit         RateLimitConfig   `json:"rate_limit" yaml:"rate_limit"`
	SSE               SSEConfig         `json:"sse" yaml:"sse"`
	LogRequests       bool              `json:"log_requests" yaml:"log_requests"`
	LogLevel          string            `json:"log_level" yaml:"log_level"`
}

// RateLimitConfig configures rate limiting for MCP requests
type RateLimitConfig struct {
	Enabled           bool `json:"enabled" yaml:"enabled"`
	RequestsPerMinute int  `json:"requests_per_minute" yaml:"requests_per_minute"`
	Burst             int  `json:"burst" yaml:"burst"`
}

// SSEConfig configures Server-Sent Events behavior
type SSEConfig struct {
	HeartbeatInterval     int `json:"heartbeat_interval" yaml:"heartbeat_interval"`         // Seconds
	ConnectionTimeout     int `json:"connection_timeout" yaml:"connection_timeout"`         // Seconds
	MaxConnectionsPerUser int `json:"max_connections_per_user" yaml:"max_connections_per_user"`
}

// DefaultConfig returns the default MCP plugin configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		EnabledOperations: []string{"crud", "schema"},
		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 60,
			Burst:             10,
		},
		SSE: SSEConfig{
			HeartbeatInterval:     30,
			ConnectionTimeout:     300,
			MaxConnectionsPerUser: 3,
		},
		LogRequests: true,
		LogLevel:    "info",
	}
}

// ParseConfig parses configuration from various formats
func ParseConfig(cfg interface{}) (*Config, error) {
	if cfg == nil {
		return DefaultConfig(), nil
	}

	// Handle map[string]interface{} from YAML/JSON parsing
	if configMap, ok := cfg.(map[string]interface{}); ok {
		config := DefaultConfig()

		if enabled, ok := configMap["enabled"].(bool); ok {
			config.Enabled = enabled
		}

		if ops, ok := configMap["enabled_operations"].([]interface{}); ok {
			config.EnabledOperations = make([]string, len(ops))
			for i, op := range ops {
				if opStr, ok := op.(string); ok {
					config.EnabledOperations[i] = opStr
				}
			}
		}

		if rateLimitMap, ok := configMap["rate_limit"].(map[string]interface{}); ok {
			if enabled, ok := rateLimitMap["enabled"].(bool); ok {
				config.RateLimit.Enabled = enabled
			}
			if rpm, ok := rateLimitMap["requests_per_minute"].(int); ok {
				config.RateLimit.RequestsPerMinute = rpm
			}
			if burst, ok := rateLimitMap["burst"].(int); ok {
				config.RateLimit.Burst = burst
			}
		}

		if sseMap, ok := configMap["sse"].(map[string]interface{}); ok {
			if heartbeat, ok := sseMap["heartbeat_interval"].(int); ok {
				config.SSE.HeartbeatInterval = heartbeat
			}
			if timeout, ok := sseMap["connection_timeout"].(int); ok {
				config.SSE.ConnectionTimeout = timeout
			}
			if maxConns, ok := sseMap["max_connections_per_user"].(int); ok {
				config.SSE.MaxConnectionsPerUser = maxConns
			}
		}

		if logReq, ok := configMap["log_requests"].(bool); ok {
			config.LogRequests = logReq
		}

		if logLevel, ok := configMap["log_level"].(string); ok {
			config.LogLevel = logLevel
		}

		return config, nil
	}

	// Handle *Config directly
	if config, ok := cfg.(*Config); ok {
		return config, nil
	}

	return nil, fmt.Errorf("unsupported config type: %T", cfg)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	// Validate enabled operations
	validOps := map[string]bool{"crud": true, "schema": true}
	for _, op := range c.EnabledOperations {
		if !validOps[op] {
			return fmt.Errorf("invalid operation: %s (valid: crud, schema)", op)
		}
	}

	// Validate rate limiting
	if c.RateLimit.Enabled {
		if c.RateLimit.RequestsPerMinute <= 0 {
			return fmt.Errorf("requests_per_minute must be > 0")
		}
		if c.RateLimit.Burst <= 0 {
			return fmt.Errorf("burst must be > 0")
		}
	}

	// Validate SSE settings
	if c.SSE.HeartbeatInterval <= 0 {
		return fmt.Errorf("heartbeat_interval must be > 0")
	}
	if c.SSE.ConnectionTimeout <= 0 {
		return fmt.Errorf("connection_timeout must be > 0")
	}
	if c.SSE.MaxConnectionsPerUser <= 0 {
		return fmt.Errorf("max_connections_per_user must be > 0")
	}

	// Validate log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid log_level: %s (valid: debug, info, warn, error)", c.LogLevel)
	}

	return nil
}

// GetHeartbeatInterval returns the heartbeat interval as time.Duration
func (c *Config) GetHeartbeatInterval() time.Duration {
	return time.Duration(c.SSE.HeartbeatInterval) * time.Second
}

// GetConnectionTimeout returns the connection timeout as time.Duration
func (c *Config) GetConnectionTimeout() time.Duration {
	return time.Duration(c.SSE.ConnectionTimeout) * time.Second
}

// IsOperationEnabled checks if a specific operation is enabled
func (c *Config) IsOperationEnabled(operation string) bool {
	for _, op := range c.EnabledOperations {
		if op == operation {
			return true
		}
	}
	return false
}
