package main

import (
	"context"
	"fmt"
	"net"

	"github.com/sanyewudezhuzi/gRPC_study/pb"

	"google.golang.org/grpc"
)

// hello_server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("hello: " + req.Name)
	return &pb.HelloResponse{Reply: "Hello " + req.Name}, nil
}

func main() {
	// 监听本地的 9090 端口
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		panic("failed to listen")
	}
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()
	// 在 grpc 客户端注册我们自己编写的服务
	pb.RegisterGreeterServer(grpcServer, &server{})
	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		panic("failed to server")
	}
}
