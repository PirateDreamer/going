package http

import (
	"context"
	"reflect"
	"time"

	"github.com/PirateDreamer/going/xerr"
	"github.com/cloudwego/hertz/pkg/app"
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

func ResSuccess(ctx context.Context, c *app.RequestContext, data any) {
	if data == nil {
		data = Empty{}
	}
	c.JSON(200, Response{
		Code:  0,
		Msg:   "success",
		Time:  time.Now().Unix(),
		ReqId: ctx.Value("traceId").(string),
		Data:  data,
	})
}

func ResFail(ctx context.Context, c *app.RequestContext, err error) {
	ResFailWithData(ctx, c, err, Empty{})
}

func ResFailWithData(ctx context.Context, c *app.RequestContext, err error, data any) {
	result := Response{
		Time:  time.Now().Unix(),
		ReqId: ctx.Value("traceId").(string),
		Data:  data,
	}
	switch e := err.(type) {
	case xerr.BizError:
		result.Code = e.Code
		result.Msg = e.Msg
		result.Err = err.Error()
	default:
		result.Code = 1
		result.Msg = "系统繁忙"
		result.Err = err.Error()
	}
	c.JSON(200, result)
}

func HertzWrapHandler[T, R any](fn func(ctx context.Context, c *app.RequestContext, req T) (*R, error)) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 参数绑定
		var req T
		t := reflect.TypeOf(req)
		if t != reflect.TypeOf(Empty{}) {
			if err := c.Bind(&req); err != nil {
				ResFail(ctx, c, xerr.NewCommBizErr(err.Error()))
				c.Abort()
				return
			}
		}
		// 执行处理逻辑
		resp, err := fn(ctx, c, req)

		// 结果处理
		if err != nil {
			ResFail(ctx, c, err)
		} else {
			ResSuccess(ctx, c, resp)
		}
		c.Next(ctx)
	}
}
