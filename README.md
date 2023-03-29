# v3.2 双向流式 rpc

1. 定义服务
    
    ```protobuf
    rpc BidiHello(stream HelloRequest) returns (stream HelloResponse) {}
    ```
    
2. 服务端实现
    
    ```go
    // 从隔壁大佬偷的价值连城的 AI 模型
    func aimodel(s string) string {
            s = strings.ReplaceAll(s, "吗", "")
            s = strings.ReplaceAll(s, "吧", "")
            s = strings.ReplaceAll(s, "你", "我")
            s = strings.ReplaceAll(s, "？", "!")
            s = strings.ReplaceAll(s, "?", "!")
            return s
    }

    // BidiHello 双向流式打招呼
    func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
            for {
                    // 接收流式请求
                    res, err := stream.Recv()
                    if res == nil {
                            return err
                    }
                    if err != nil {
                            log.Fatalln("failed to recv:", err)
                            return err
                    }
                    // 对收到的数据做些处理
                    fmt.Println(res.GetName())
                    reply := aimodel(res.GetName())
                    // 返回流式响应
                    if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
                            return err
                    }
            }
    }
    ```
    
3. 客户端实现
    
    ```go
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
    ```
    
4. 运行

    <div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/35c5fe82a9794b7f95c3b7129758f003~tplv-k3u1fbpfcp-watermark.image?" /></div>
