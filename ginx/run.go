package ginx

import (
	"context"
	"time"

	"github.com/PirateDreamer/going/comm"
	"github.com/PirateDreamer/going/grpcx"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type Response struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Err   string `json:"err"`
	Data  any    `json:"data"`
	ReqId string `json:"req_id"`
	Time  int64  `json:"time"`
}

type Empty struct{}

func ResSuccess(ctx context.Context, c *gin.Context, data any) {
	if data == nil {
		data = Empty{}
	}

	c.JSON(200, Response{
		Code:  0,
		Msg:   "success",
		Time:  time.Now().Unix(),
		ReqId: comm.GetReqId(ctx),
		Data:  data,
	})
}

func ResFail(ctx context.Context, c *gin.Context, err error) {
	ResFailWithData(ctx, c, err, Empty{})
}

func ResFailWithData(ctx context.Context, c *gin.Context, err error, data any) {
	result := Response{
		Time:  time.Now().Unix(),
		Data:  data,
		ReqId: comm.GetReqId(ctx),
	}
	st, ok := status.FromError(err)
	if ok && st.Code() >= 1000 {
		// 业务错误
		result.Code = int(st.Code())
		result.Msg = st.Message()
		result.Err = st.Err().Error()
	} else {
		result.Code = 10000
		result.Msg = "系统繁忙"
		result.Err = err.Error()
	}
	c.JSON(200, result)
}

func Run[T, R any](fn func(ctx context.Context, c *gin.Context, req T) (*R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 参数绑定
		var req T
		if err := c.ShouldBind(&req); err != nil {
			ResFail(ctx, c, grpcx.BizErr(err.Error()))
			c.Abort()
			return
		}

		// 执行处理逻辑
		resp, err := fn(ctx, c, req)

		// 结果处理
		if err != nil {
			ResFail(ctx, c, err)
		} else {
			ResSuccess(ctx, c, resp)
		}
		c.Next()
	}
}
