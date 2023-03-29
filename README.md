# v1 入门示例

1. 在 server 端创建 *pb* 文件夹，编写 *hello.proto* 文件来定义服务

    ```protobuf
    // 版本声明，使用Protocol Buffers v3版本
    syntax = "proto3";

    // 这部分的内容是关于最后生成的 go 文件是处在哪个目录哪个包中
    // . 代表在当前目录生成
    // 以 ; 号间隔
    // service 代表了生成的 go 文件的包名是 service 
    option go_package = ".;pb";

    // 然后我们需要定义一个服务，在这个服务中需要有一个方法，这个方法可以接受客户端的参数，再返回服务端的响应
    // 其实很容易可以看出，我们定义了一个 service ，称为 Greeter ，这个服务中有一个 rpc 方法，名为 SayHello 
    // 这个方法会发送一个 HelloRequest ，然后返回一个 HelloResponse 
    service Greeter {
        rpc SayHello (HelloRequest) returns (HelloResponse) {}
    }

    // message 关键字，其实你可以理解为 Golang 中的结构体
    // 这里比较特别的是变量后面的“赋值”
    // 这里并不是赋值，而是定义在这个 message 中的位置
    message HelloRequest {
        string name = 1;    // 第一行就标识为 1
        // int64 age = 2;   // 第二行就标识为 2
    }

    message HelloResponse {
        string reply = 1;
    }
    ```

2. 在 client 端也创建一个 *pb* 文件夹，并将上面的 *proto* 文件拷贝到客户端的文件夹中
3. 分别在服务端和客户端的 *pb* 目录下打开命令行，并执行下面两条命令：
    * `protoc --go_out=. hello.proto`
    * `protoc --go-grpc_out=. hello.proto`
    
   *pb* 目录中会自动生成两个文件：
    * hello.pb.go
    * hello_grpc.pb.go
4. 此时服务端的目录结构：<div align=center><img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b3b4a45e66674fbca24a80e905c4d121~tplv-k3u1fbpfcp-watermark.image?" /></div>

   以及客户端的目录结构：<div align=center><img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a8bc01cfdcc74386982d17b44438c044~tplv-k3u1fbpfcp-watermark.image?" /></div>

5. 编写服务端 *main.go* 文件

    ```go
    package main

    import (
            "context"
            "fmt"
            "net"

            "github.com/sanyewudezhuzi/gRPC_study/pb"

            "google.golang.org/grpc"
    )

    // hello_server

    type server struct {
            pb.UnimplementedGreeterServer
    }

    func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
            fmt.Println("hello: " + req.Name)
            return &pb.HelloResponse{Reply: "Hello " + req.Name}, nil
    }

    func main() {
            // 监听本地的 9090 端口
            listen, err := net.Listen("tcp", "127.0.0.1:9090")
            if err != nil {
                    panic("failed to listen")
            }
            // 创建gRPC服务器
            grpcServer := grpc.NewServer()
            // 在 grpc 客户端注册我们自己编写的服务
            pb.RegisterGreeterServer(grpcServer, &server{})
            // 启动服务
            err = grpcServer.Serve(listen)
            if err != nil {
                    panic("failed to server")
            }
    }
    ```
    
6. 编写服务端 *main.go* 文件

    ```go
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
    ```

7. 先后执行服务端和客户端：<div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/0a8d2c698e0744458d57df6a21c49525~tplv-k3u1fbpfcp-watermark.image?" /></div>
