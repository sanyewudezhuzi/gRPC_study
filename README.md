# v3.0 服务端流式 rpc

1. 定义服务
    
    ```protobuf
    rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse) {}
    ```
    
2. 服务端实现
    
    ```go
    func (s *server) LotsOfReplies(req *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
            fmt.Println(req.GetName())
            words := []string{
                    "你好 ",
                    "hello ",
                    "こんにちは ",
                    "안녕하세요 ",
                    "สวัสดี ",
            }
            for _, word := range words {
                    data := &pb.HelloResponse{
                            Reply: word + req.GetName(),
                    }
                    // 使用 send 方法返回多个数据
                    if err := stream.Send(data); err != nil {
                            return err
                    }
            }
            return nil
    }
    ```
    
3. 客户端实现
    
    ```go
    // LotsOfReplies 返回使用多种语言打招呼
    func runLotsOfReplies(c pb.GreeterClient) {
            // server端流式RPC
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()
            stream, err := c.LotsOfReplies(ctx, &pb.HelloRequest{Name: *name})
            if err != nil {
                    log.Fatalln("failed to lotsofreplies:", err)
            }
            for {
                    // 接收服务端返回的流式数据，当收到io.EOF或错误时退出
                    res, err := stream.Recv()
                    if err == io.EOF {
                            break
                    }
                    if err != nil {
                            log.Fatalln("failed to recv:", err)
                    }
                    fmt.Println(res.GetReply())
            }
    }
    ```
    
4. 运行

    <div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e53667d6030d4528b9f6d9210e2634e7~tplv-k3u1fbpfcp-watermark.image?" /></div>
