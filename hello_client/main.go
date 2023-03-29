package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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
)

func runBidiHello(c pb.GreeterClient) {
	// 双向流模式
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	stream, err := c.BidiHello(ctx)
	// 双向流要手动关闭来标记消息的结束
	defer stream.CloseSend()
	if err != nil {
		log.Fatalln("failed to bidihello:", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 接收服务端返回的响应
			res, err := stream.Recv()
			if res == nil {
				return
			}
			if err != nil {
				log.Fatalln("failed to recv:", err)
			}
			fmt.Println("AI: ", res.GetReply())
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
	// 客户端调用 BidiHello 方法，一边从终端获取输入的请求数据发送至服务端，一边从服务端接收流式响应
	runBidiHello(client)
}
