package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sanyewudezhuzi/gRPC_study/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// hello_server

type server struct {
	pb.UnimplementedGreeterServer
}

var database = map[string]string{
	"zhuzi": "abc123",
}

// normalCall 普通调用
func (s *server) normalCall(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 设置 trailer
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()
	// 获取 metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	name := md["name"][0]
	pwd := md["pwd"][0]
	if p, ok := database[name]; ok && p == pwd {
		fmt.Println("name:", name)
	} else {
		return nil, status.Error(codes.Unauthenticated, "invalid info")
	}
	// 发送 header
	header := metadata.New(map[string]string{"greeting": name + " say: hello " + req.Name})
	grpc.SendHeader(ctx, header)
	return &pb.HelloResponse{Reply: req.Name}, nil
}

// streamCalls 流式调用
func (s server) streamCalls(stream pb.Greeter_BidiHelloServer) error {
	// 设置 trailer
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()
	// 获取 metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	name := md["name"][0]
	pwd := md["pwd"][0]
	if p, ok := database[name]; ok && p == pwd {
		fmt.Println("name:", name)
	} else {
		return status.Error(codes.Unauthenticated, "invalid info")
	}
	// 发送数据
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalln("failed to recv:", err)
			return err
		}
		fmt.Println("get name: ", res.GetName())
		if err := stream.Send(&pb.HelloResponse{Reply: name + "say: hello " + res.Name}); err != nil {
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
