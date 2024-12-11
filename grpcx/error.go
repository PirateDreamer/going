package grpcx

import "google.golang.org/grpc/status"

// BizError 业务错误
func BizErr(message string, a ...any) error {
	return status.Errorf(1000, message, a...)
}
