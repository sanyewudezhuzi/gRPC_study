package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sanyewudezhuzi/gRPC_study/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// hello_client

const (
	defaultName = "world"
)

var (
	name = flag.String("name", defaultName, "name")
	addr = flag.String("addr", "127.0.0.1:9090", "IP address and port number of the tcp connection")
)

func runSayHello(c pb.GreeterClient, name string) {
	// 创建 metadata
	md := metadata.Pairs(
		"name", "zhuzi",
		"pwd", "abc123",
	)
	// 基于 metadata 创建 context
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// RPC 调用
	var header, trailer metadata.MD
	res, err := c.SayHello(
		ctx,
		&pb.HelloRequest{Name: name},
		grpc.Header(&header),   // 接收服务端发来的 header
		grpc.Trailer(&trailer), // 接收服务端发来的 trailer
	)
	if err != nil {
		log.Println("failed to call sayhello:", err)
		return
	}
	// 从 header 中获取 greeting
	if g, ok := header["greeting"]; ok {
		fmt.Println(g[0])
	} else {
		log.Println("failed to get greeting")
		return
	}
	// 获取响应
	fmt.Println("res: ", res.GetReply())
	// 从trailer中取timestamp
	if t, ok := trailer["timestamp"]; ok {
		fmt.Println("timestamp: ", t[0])
	} else {
		log.Println("failed to get timestamp")
	}
}

func runBidiHello(c pb.GreeterClient) {
	// 创建 metadata
	md := metadata.Pairs(
		"name", "zhuzi",
		"pwd", "abc123",
	)
	// 基于 metadata 创建 context
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// RPC 调用
	stream, err := c.BidiHello(ctx)
	defer stream.CloseSend()
	if err != nil {
		log.Println("failed to call bidhello:", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 接收服务端返回的响应
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalln("failed to recv:", err)
			}
			fmt.Println(res.GetReply())
		}
	}()
	// 从标准输入获取用户输入
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "EXIT" {
			stream.Send(nil)
			break
		}
		// 将获取到的数据发送至服务端
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err != nil {
			log.Fatalln("failed to send:", err)
		}
	}
	wg.Wait()
	// 结束时读取 trailer
	trailer := stream.Trailer()
	if t, ok := trailer["timestamp"]; ok {
		fmt.Println("timestamp: ", t[0])
	} else {
		log.Println("failed to get timestamp")
	}
}

func main() {
	// 将用户传递至命令行的参数解析为对应变量值
	flag.Parse()
	// 连接到server端
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to connect:", err)
	}
	defer conn.Close()
	// 建立连接
	client := pb.NewGreeterClient(conn)
	// 普通调用
	// runSayHello(client, *name)
	// 流式调用
	runBidiHello(client)
}
