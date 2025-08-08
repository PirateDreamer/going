package ginc

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Empty struct{}

func Run[T, R any](fn func(context.Context, *gin.Context, T) (*R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := ctx.Request.Context()

		// 参数绑定
		var req T
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// 执行函数
		resp, err := fn(c, ctx, req)

		// 结果处理
		resHandler(c, ctx)
		if err != nil {
			ResFailWithData(c, ctx, err, resp)
		} else {
			ResSuccess(c, ctx, resp)
		}
		ctx.Next()
	}
}
