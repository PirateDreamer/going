package sidecar

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

// grpcProxy 是一个通用的 gRPC 代理
type grpcProxy struct {
	backendAddr string // 后端 gRPC 服务器地址
}

// ProxyStream 是一个通用的流式代理方法
func (p *grpcProxy) ProxyStream(srv interface{}, stream grpc.ServerStream) error {
	// 连接到后端 gRPC 服务器
	conn, err := grpc.Dial(p.backendAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// 获取客户端上下文
	ctx := stream.Context()

	// 从上下文中提取元数据
	md, _ := metadata.FromIncomingContext(ctx)

	method, _ := grpc.MethodFromServerStream(stream)

	// 创建到后端服务器的流
	clientStream, err := grpc.NewClientStream(ctx, &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, conn, method, grpc.Header(&md))
	if err != nil {
		return err
	}

	// 在 goroutine 中从客户端接收数据并转发到后端
	go func() {
		for {
			var msg []byte
			if err := stream.RecvMsg(&msg); err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Failed to receive message from client: %v", err)
				return
			}
			if err := clientStream.SendMsg(msg); err != nil {
				log.Printf("Failed to send message to backend: %v", err)
				return
			}
		}
	}()

	// 从后端接收数据并转发到客户端
	for {
		var msg []byte
		if err := clientStream.RecvMsg(&msg); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Failed to receive message from backend: %v", err)
			return err
		}
		if err := stream.SendMsg(msg); err != nil {
			log.Printf("Failed to send message to client: %v", err)
			return err
		}
	}

	return nil
}
