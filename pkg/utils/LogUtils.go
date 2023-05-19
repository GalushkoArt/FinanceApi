package utils

import (
	"context"
	"github.com/rs/zerolog"
)

func GetRequestId(ctx context.Context) string {
	if id, ok := ctx.Value("requestid").(string); ok {
		return id
	}
	return ""
}

func LogRequest(c context.Context, e *zerolog.Event) *zerolog.Event {
	if e == nil {
		return nil
	}
	if id, ok := c.Value("requestid").(string); ok {
		return e.Str("requestid", id)
	}
	return e
}
