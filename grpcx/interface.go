package grpcx

import (
	"context"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 检查并创建请求ID
func GrpcCheckReqId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqId := ctx.Value("req_id")
		if reqId == nil {
			traceId := uuid.NewV4().String()
			ctx = context.WithValue(ctx, "req_id", traceId)
		}
		return handler(ctx, req)
	}
}

func GrpcRequestLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

// 参数校验
func ValidationInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// 检查请求是否有 Validate() 方法
	if validator, ok := req.(interface{ Validate() error }); ok {
		// 如果有 Validate 方法，调用它进行参数校验
		if err := validator.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}
	} else {
		// 如果请求不包含 Validate 方法，可能是某些特殊消息类型
		return nil, status.Errorf(codes.InvalidArgument, "Request type %v does not support validation", reflect.TypeOf(req))
	}

	// 如果校验通过，继续调用下一个 handler（实际的 gRPC 方法）
	return handler(ctx, req)
}
