# v4 metadata

在多个微服务的调用当中，信息交换常常是使用方法之间的参数传递的方式，但是在有些场景下，一些信息可能和 RPC 方法的业务参数没有直接的关联，所以不能作为参数的一部分，在 gRPC 中，可以使用元数据来存储这类信息。元数据最主要的作用：
1.  提供RPC调用的元数据信息，例如用于链路追踪的traceId、调用时间、应用版本等等。
2.  控制gRPC消息的格式，例如是否压缩或是否加密。

实际开发中，我们一般使用第三方库 [google.golang.org/grpc/metadata](https://pkg.go.dev/google.golang.org/grpc/metadata) 来操作元数据。

## 创建元数据

元数据的类型定义：

```go
type md map[string][]string
```

metadata中的键是大小写不敏感的，由字母、数字和特殊字符 `-` 、`_` 、`.` 组成并且不能以 `grpc-` 开头（gRPC保留自用），二进制值的键名必须以 `-bin` 结尾

1. 使用 metadata 库的 `New()` 函数来创建元数据
    
    ```go
    md := metadata.New(map[string]string{"key1": "val1", "key2": "val2"})
    ```
2. 使用 metadata 库的 `Pairs()` 函数来创建元数据

    ```go
    md := metadata.Pairs(
        "key1", "val1",
        "key1", "val1-2", // "key1"的值将会是 []string{"val1", "val1-2"}
        "key2", "val2",
        "key-bin", string([]byte{96, 102}), // 二进制数据在发送前会进行(base64) 编码
    )
    ```
    
3. 使用 `FromIncomingContext` 从 RPC 请求的上下文中获取

    ```go
    func (s *server) SomeRPC(ctx context.Context, in *pb.SomeRequest) (*pb.SomeResponse, err) {
        md, ok := metadata.FromIncomingContext(ctx)
        // do something with metadata
    }
    ```

## 客户端发送接收元数据

```go
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
```

## 服务端发送接收元数据

```go
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
```
