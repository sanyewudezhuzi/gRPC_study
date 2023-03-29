package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
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
	addr  = flag.String("addr", "127.0.0.1:9090", "IP address and port number of the tcp connection")
	names = flag.String("names", defaultName, "Names to greet")
)

func runLotsOfGreeting(c pb.GreeterClient) {
	// 客户端流式RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalln("failed to lotsofgreetings:", err)
	}
	namelist := strings.Split(*names, ",")
	for _, name := range namelist {
		// 发送流式数据
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalln("failed to send:", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("failed to closeandrecv:", err)
	}
	fmt.Println(res.GetReply())
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
	// 客户端调用 LotsOfGreetings 方法，向服务端发送流式请求数据，接收返回值并打印
	runLotsOfGreeting(client)
}
