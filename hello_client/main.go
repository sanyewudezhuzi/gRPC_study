package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/sanyewudezhuzi/gRPC_study/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// hello_client

const (
	defaultName = "world"
	cretFile    = "E:\\GoLanguage\\src\\gRPC" + "\\gRPC_study\\hello_server\\key\\test.pem" // 证书文件
	defaultDNS  = "*.sanyewu.com"                                                           // openssl.cfg 文件中 [ alt_names ] 标签定义的 DSN
)

var (
	addr               = flag.String("addr", "127.0.0.1:9090", "IP address and port number of the tcp connection")
	name               = flag.String("name", defaultName, "Name to greet")
	serverNameOverride = flag.String("dns", defaultDNS, "Server name override")
)

func main() {
	// 将用户传递至命令行的参数解析为对应变量值
	flag.Parse()
	// 获取证书
	creds, err := credentials.NewClientTLSFromFile(cretFile, *serverNameOverride)
	if err != nil {
		panic("failed to creds")
	}
	// 连接到server端
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
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
