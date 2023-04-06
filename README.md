# gRPC_study
简单介绍了 gRPC 的使用

# 前言
* 真聊会天
    * 本文开袋即食，大致内容目录里面可以体现，主要帮助新手快速学习 GRPC 

* 开发环境
    * Windows x64
    * go version go1.20 windows/amd64
    * code --version 1.76.2
    * libprotoc 22.1
    * OpenSSL 3.1.0 14 Mar 2023 (Library: OpenSSL 3.1.0 14 Mar 2023)

* 代码地址
    * https://github.com/sanyewudezhuzi/gRPC_study

---

# v0 准备工作

1. 根据自己的开发环境下载 [*protoc*](https://github.com/protocolbuffers/protobuf/releases) ，并将其 *bin* 目录配置到环境变量中
    * 在命令行中输入 `protoc --version` 即可查看是否配置成功
    * ~~如果你懒得搜且开发环境和我类似，你也可以和我[下载一样的](https://github.com/protocolbuffers/protobuf/releases/download/v22.1/protoc-22.1-win64.zip)~~
2. 创建项目与 go mod 文件，在命令行中输入 `go get google.golang.org/grpc` 安装 grpc 的核心库
3. 在命令行中输入 `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` 和 `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` 安装对应工具
    * 在命令行中分别输入 `protoc-gen-go --version` 和 `protoc-gen-go-grpc --version` 即可查看是否配置成功
    * 此时在 *%GOPATH%/bin* 目录下可以看到这两个工具：<div align=center><img src="https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/dbd018e7d5f240af91b7077b672291eb~tplv-k3u1fbpfcp-watermark.image?" /></div>
