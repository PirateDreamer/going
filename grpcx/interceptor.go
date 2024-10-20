package grpcx

import (
	"context"
	"log"
	"reflect"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 打印 panic 信息和堆栈
			log.Printf("Panic recovered: %v\n%s", r, debug.Stack())
			// 返回 gRPC 内部错误
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()

	// 调用实际的 gRPC handler
	return handler(ctx, req)
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
