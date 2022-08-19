package internal

import (
	"context"

	"github.com/google/uuid"
)

func RequestId(ctx context.Context) string {
	if id, ok := ctx.Value("x-request-id").(string); ok {
		return id
	}

	return uuid.NewString()
}
