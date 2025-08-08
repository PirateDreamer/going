package ginc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PirateDreamer/going/comm"
	"github.com/PirateDreamer/going/xerr"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Err  string `json:"err"`
	Data any    `json:"data"`
}

func ResFail(ctx context.Context, c *gin.Context, err error) {
	ResFailWithData(ctx, c, err, nil)
}

func ResFailWithData(ctx context.Context, c *gin.Context, err error, data any) {
	// 判断错误类型,并转化错误
	var resultErr xerr.BizError
	switch e := err.(type) {
	case *xerr.BizError:
		resultErr = *e
	default:
		resultErr = xerr.BizError{
			Code: "1",
			Msg:  "服务不见了",
		}
	}
	res := Response{Code: resultErr.Code, Msg: resultErr.Msg, Err: err.Error()}
	if data != nil {
		res.Data = data
	} else {
		res.Data = Empty{}
	}
	c.JSON(http.StatusOK, res)
}

func ResSuccess(ctx context.Context, c *gin.Context, data any) {
	res := Response{Code: "0", Msg: "success"}
	if data == nil {
		res.Data = Empty{}
	} else {
		res.Data = data
	}
	c.JSON(http.StatusOK, res)
}

func resHandler(ctx context.Context, c *gin.Context) {
	c.Request.Response.Header.Set("X-Trace-ID", comm.GetTraceID(ctx))
	c.Request.Response.Header.Set("X-Response-Time", fmt.Sprintf("%d", time.Now().UnixMicro()))
}
