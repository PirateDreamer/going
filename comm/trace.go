package comm

import (
	"context"

	"github.com/google/uuid"
)

func GetTraceID(c context.Context) string {
	value := c.Value("trace_id")
	if value != nil || value.(string) != "" {
		return value.(string)
	}
	traceID := uuid.New().String()
	c = context.WithValue(c, "trace_id", traceID)
	return traceID
}
