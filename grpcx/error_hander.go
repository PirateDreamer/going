package grpcx

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// 始终返回 200
	w.WriteHeader(http.StatusOK)
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
