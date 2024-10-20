package grpcx

import (
	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/gormx"
	"github.com/PirateDreamer/going/gredis"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Server struct {
	GrpcServer *GRPCServer
	HTTPServer *HTTPServer
}

type ServerConfig struct {
	GrpcInterceptors []grpc.UnaryServerInterceptor
}

func NewServer(args ...any) *Server {
	conf.InitConfig(nil)
	gormx.InitMysql()
	gredis.InitRedis()

	// 服务配置
	var serverConfig ServerConfig
	if len(args) > 0 {
		serverConfig = args[0].(ServerConfig)
	}

	server := Server{}

	server.GrpcServer = NewGrpcServer(serverConfig.GrpcInterceptors)
	server.HTTPServer = NewHttpServer()

	return &server
}

// Start 启动服务
func (s *Server) Start() {
	go s.StartGrpc()
	go s.StartHttp()

	select {}
}

func (server *Server) StartGrpc() {
	server.GrpcServer.StartGrpcServer(viper.GetString("grpc.addr"))
}

func (server *Server) StartHttp() {
	server.HTTPServer.StartHttpServer(viper.GetString("http.addr"))
}
