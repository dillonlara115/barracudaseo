package api

import (
	"context"
)

type contextKey string

const userIDKey contextKey = "user_id"

// contextWithUserID adds user ID to context
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// userIDFromContext extracts user ID from context
func userIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

