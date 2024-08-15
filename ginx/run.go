package ginx

import (
	"reflect"
	"time"

	"github.com/PirateDreamer/going/comm"
	"github.com/PirateDreamer/going/xerr"
	"github.com/gin-gonic/gin"
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

func ResSuccess(c *gin.Context, data any) {
	if data == nil {
		data = Empty{}
	}

	c.JSON(200, Response{
		Code:  0,
		Msg:   "success",
		Time:  time.Now().Unix(),
		ReqId: comm.GetReqId(c),
		Data:  data,
	})
}

func ResFail(c *gin.Context, err error) {
	ResFailWithData(c, err, Empty{})
}

func ResFailWithData(c *gin.Context, err error, data any) {
	result := Response{
		Time:  time.Now().Unix(),
		Data:  data,
		ReqId: comm.GetReqId(c),
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

func Run[T, R any](fn func(c *gin.Context, req T) (*R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 参数绑定
		var req T
		t := reflect.TypeOf(req)
		if t != reflect.TypeOf(Empty{}) {
			if err := c.Bind(&req); err != nil {
				ResFail(c, xerr.NewCommBizErr(err.Error()))
				c.Abort()
				return
			}
		}
		// 执行处理逻辑
		resp, err := fn(c, req)

		// 结果处理
		if err != nil {
			ResFail(c, err)
		} else {
			ResSuccess(c, resp)
		}
		c.Next()
	}
}
