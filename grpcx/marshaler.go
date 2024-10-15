package grpcx

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type CustomMarshaler struct {
	runtime.Marshaler
}

// Marshal 实现自定义的响应格式
func (c *CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	// 是否为错误码
	if res, ok := v.(*status.Status); ok {
		errorResponse := map[string]interface{}{
			"code": res.Code,
			"msg":  res.Message,
			"data": res.Details,
		}

		return c.Marshaler.Marshal(errorResponse)
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": v,
	}
	return c.Marshaler.Marshal(response)
}
