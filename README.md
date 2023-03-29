# v2 加密认证

先唠嗑几个知识点：

> key：服务器上的私钥文件，用于对发送给客户端数据的加密，以及对客户端接收到数据的解密
>
> csr：证书签名请求文件，用于提交给证书颁发机构（CA）对证书签名
>
> crt：有证书颁发机构（CA）签名后的证书，或者是开发者自签名的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息
>
> pem：是基于 Base64 编码的证书格式，扩展名包括 PEM 、CRT 、CER

在 v1 的基础上，于服务端中**新增 *key* 目录**

1. 下载 [openssl](http://slproweb.com/products/Win32OpenSSL.html)
    * 在命令行中输入 `openssl version` 即可查看是否配置成功
    * ~~如果你懒得搜且开发环境和我类似，你也可以和我[下载一样的](https://slproweb.com/download/Win64OpenSSL-3_1_0.exe)~~
2. 在 *key* 目录下打开命令行，输入 `openssl genrsa -out server.key 2048` 生成私钥
    * 跟我一起念~ 私(sī)钥(yuè)~
3. 在命令行输入 `openssl req -new -x509 -key server.key -out server.crt -days 36500` 生成证书，有关选项可以不填，全部回车即可
4. 在命令行输入 `openssl req -new -key server.key -out server.csr` 生成 csr ，有关选项可以不填，全部回车即可
5. 更改 *openssl.cfg*
    * 在 openssl 安装路径的 *bin* 目录下找到 *openssl.cfg* 文件，将其拷贝到项目的 *key* 目录下
    * 搜索 `copy_extensions` ，将 `copy_extensions = copy` 开启<div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2bed0b8985384ab192dc20b4f2411294~tplv-k3u1fbpfcp-watermark.image?" /></div>
    * 搜索 `req_extensions` ，将 `req_extensions = v3_req # The extensions to add to a certificate request` 开启<div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f064a6b62f43437c97c37081b2bdcebf~tplv-k3u1fbpfcp-watermark.image?" /></div>
    * 搜索 `v3_req` ，添加字段 `**subjectAltName** = @alt_names`<div align=center><img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a0776e7585ad4a508d7af524bdb0de0b~tplv-k3u1fbpfcp-watermark.image?" /></div>
    * 添加标签 `[ alt_names ]` 和字段 `DNS.1 = *.sanyewu.com`<div align=center><img src="https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/dc7b0ed77d36416b9d1ceea71dd9789e~tplv-k3u1fbpfcp-watermark.image?" /></div>
      其中 sanyewu 可以替换为你所想要的域名
7. 在命令行输入 `openssl genpkey -algorithm RSA -out test.key` 生成证书私钥
8. 在命令行输入 `openssl req -new -nodes -key test.key -out test.csr -days 3650 -subj "/C=cn/OU=myorg/O=mycomp/CN=myname" -config ./openssl.cfg -extensions v3_req` 通过私钥 test.key 生成证书请求文件 test.csr 
9. 在命令行输入 `openssl x509 -req -days 365 -in test.csr -out test.pem -CA server.crt -CAkey server.key -CAcreateserial -extfile ./openssl.cfg -extensions v3_req` 生成 SAN 证书 pem 
10. 完成后你的 key 目录下应该有这 8 个文件：<div align=center><img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/26be4295142e47498bd05f6b960976de~tplv-k3u1fbpfcp-watermark.image?" /></div>

**更新服务端的 *main.go***

```go
func main() {
	// 生成证书
        // cretFile: 证书文件的绝对路径
        // keyFile: 私钥文件的绝对路径
	creds, err := credentials.NewServerTLSFromFile(cretFile, keyFile)
	if err != nil {
		panic("failed to creds")
	}
	// 监听本地的 9090 端口
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		panic("failed to listen")
	}
	// 创建gRPC服务器
	// 配置证书
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 在 grpc 客户端注册我们自己编写的服务
	pb.RegisterGreeterServer(grpcServer, &server{})
	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		panic("failed to server")
	}
}
```

将 *test.pem* 拷贝到客户端的 *key* 目录下，并**更新客户端的 *main.go***

```go
func main() {
	// 将用户传递至命令行的参数解析为对应变量值
	flag.Parse()
	// 获取证书
        // cretFile: 证书文件的绝对路径
        // serverNameOverride: 服务器的名称替代
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
```

分别运行服务端和客户端

<div align=center><img src="https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8181424fa23f40198ae602ac2e146a26~tplv-k3u1fbpfcp-watermark.image?" /></div>
