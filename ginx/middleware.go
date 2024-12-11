package ginx

import (
	"context"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("X-Request-ID")
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "reqId", requestId))
		c.Next()
	}
}
