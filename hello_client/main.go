package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/sanyewudezhuzi/gRPC_study/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// hello_client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:9090", "IP address and port number of the tcp connection")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	// 将用户传递至命令行的参数解析为对应变量值
	flag.Parse()
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	// 建立连接
	client := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据（这个方法在服务端实现并返回结果）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("failed to sayhello: %v", err)
	}
	fmt.Println(res.GetReply())
}
