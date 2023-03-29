# v3.1 客户端流式 rpc

1. 定义服务
    
    ```protobuf
    rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse) {}
    ```
    
2. 服务端实现
    
    ```go
    // LotsOfGreetings 接收流式数据
    func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
            reply := "你好: "
            for {
                    // 接收客户端发来的流式数据
                    res, err := stream.Recv()
                    if err == io.EOF {
                            // 最终统一回复
                            return stream.SendAndClose(&pb.HelloResponse{
                                    Reply: reply,
                            })
                    }
                    if err != nil {
                            return err
                    }
                    reply += res.GetName()
                    fmt.Println("res: ", res.GetName())
            }
    }
    ```
    
3. 客户端实现
    
    ```go
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
    ```
    
4. 运行

    <div align=center><img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/0906288183154449a1e5eface33e320b~tplv-k3u1fbpfcp-watermark.image?" /></div>
