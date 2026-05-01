package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// ContextKeyUserID stores the authenticated user ID in context
	ContextKeyUserID ContextKey = "user_id"

	// ContextKeyUserEmail stores the authenticated user email in context
	ContextKeyUserEmail ContextKey = "user_email"

	// ContextKeyUserRoles stores the authenticated user roles in context
	ContextKeyUserRoles ContextKey = "user_roles"
)

// GetUserIDFromContext extracts the user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}

// GetUserEmailFromContext extracts the user email from context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(ContextKeyUserEmail).(string)
	if !ok || email == "" {
		return "", ErrUnauthorized
	}
	return email, nil
}

// GetUserRolesFromContext extracts the user roles from context
func GetUserRolesFromContext(ctx context.Context) ([]string, error) {
	roles, ok := ctx.Value(ContextKeyUserRoles).([]string)
	if !ok {
		return nil, ErrUnauthorized
	}
	return roles, nil
}

// ToJSON converts a value to JSON string
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

// FromJSON parses JSON string into a value
func FromJSON(data string, v interface{}) error {
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// Contains checks if a slice contains a value
func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Ptr returns a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}
