package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sanyewudezhuzi/gRPC_study/pb"
	"google.golang.org/grpc"
)

// hello_server

type server struct {
	pb.UnimplementedGreeterServer
}

// 从隔壁大佬偷的价值连城的 AI 模型
func aimodel(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

// BidiHello 双向流式打招呼
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		res, err := stream.Recv()
		if res == nil {
			return err
		}
		if err != nil {
			log.Fatalln("failed to recv:", err)
			return err
		}
		// 对收到的数据做些处理
		fmt.Println(res.GetName())
		reply := aimodel(res.GetName())
		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
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
