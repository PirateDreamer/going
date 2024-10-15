package grpcx

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPServer struct {
	mux *runtime.ServeMux
}

func NewHttpServer() *HTTPServer {
	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &CustomMarshaler{
			Marshaler: &runtime.JSONPb{},
		}),
		runtime.WithErrorHandler(CustomErrorHandler),
	)
	httpServer := &HTTPServer{
		mux: mux,
	}

	return httpServer
}

func (server *HTTPServer) RegisterHandler(registerHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := registerHandlerFromEndpoint(context.Background(), server.mux, viper.GetString("grpc.addr"), opts)
	if err != nil {
		log.Fatalln("Failed to register HTTP handler:", err)
	}
}

func (server *HTTPServer) StartHttpServer(addr string) {
	log.Printf("Starting HTTP server on %s", addr)
	err := http.ListenAndServe(addr, server.mux)
	if err != nil {
		log.Fatalln("Failed to start HTTP server:", err)
	}
}
