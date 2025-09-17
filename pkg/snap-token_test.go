package connectpermify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSnapToken(t *testing.T) {
	t.Run("returns empty string when context has no snap token", func(t *testing.T) {
		ctx := context.Background()
		result := GetSnapToken(ctx)
		assert.Equal(t, "", result)
	})

	t.Run("returns snap token when present in context", func(t *testing.T) {
		expectedToken := "snap_token_123"
		ctx := context.WithValue(context.Background(), PermifySnapToken, expectedToken)
		result := GetSnapToken(ctx)
		assert.Equal(t, expectedToken, result)
	})

	t.Run("returns empty string when context value is not a string", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), PermifySnapToken, 123)
		result := GetSnapToken(ctx)
		assert.Equal(t, "", result)
	})
}

func TestSetSnapToken(t *testing.T) {
	t.Run("sets snap token in context", func(t *testing.T) {
		ctx := context.Background()
		expectedToken := "snap_token_456"

		newCtx := SetSnapToken(ctx, expectedToken)
		result := GetSnapToken(newCtx)

		assert.Equal(t, expectedToken, result)
	})

	t.Run("overwrites existing snap token", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), PermifySnapToken, "old_token")
		newToken := "new_token"

		newCtx := SetSnapToken(ctx, newToken)
		result := GetSnapToken(newCtx)

		assert.Equal(t, newToken, result)
	})
}
