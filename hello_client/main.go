package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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

func runLotsOfReplies(c pb.GreeterClient) {
	// server端流式RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v", err)
	}
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v", err)
		}
		fmt.Println(res.GetReply())
	}
}

func main() {
	// 将用户传递至命令行的参数解析为对应变量值
	flag.Parse()
	// 连接到server端
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	// 建立连接
	client := pb.NewGreeterClient(conn)
	// 客户端调用 LotsOfReplies 并将收到的数据依次打印出来
	runLotsOfReplies(client)
}
