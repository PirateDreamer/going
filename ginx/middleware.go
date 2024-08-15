package ginx

import (
	"context"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "req_id", uuid.NewV4().String()))
		c.Next()
	}
}
