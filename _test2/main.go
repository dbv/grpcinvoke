package main

import (
	"context"
	"fmt"
	"github.com/dbv/grpcinvoke"
	pb "github.com/dbv/grpcinvoke/_test2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type serverHello struct{}

func (me *serverHello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " hello!"}, nil
}

func runGprcServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := grpc.NewServer() //创建gRPC服务
	pb.RegisterHelloServer(s, &serverHello{})
	reflection.Register(s)
	// 将监听交给gRPC服务处理
	lis, err := net.Listen("tcp", ":8888") //监听所有网卡8028端口的TCP连接
	if err != nil {
		log.Fatalf("监听失败: %v", err)
		return err
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}

func main() {

	go func() {
		//启动grpc服务
		if err := runGprcServer(); err != nil {
			panic(err.Error())
		}
	}()
	time.Sleep(time.Second * 10)
	var src grpcinvoke.ProtoFile
	//导入proto common
	src.IncludeProto("./proto")
	//注册proto  定义service_hello别名，为hello.proto
	src.AddProto("service_hello", "hello.proto")
	//初始化proto 环境
	app := grpcinvoke.NewGrpcInvoke(&src, "127.0.0.1:8888", grpc.WithInsecure())

	requests := []map[string]interface{}{
		map[string]interface{}{
			"name": "luv",
		},
		map[string]interface{}{
			"name": "xiaohui",
		},
	}
	for _, request := range requests {
		output := make(map[string]interface{})
		output["message"] = nil //这里不能丢
		//proto 调用  servicename= proto中包名.proto中servicename
		if err := app.Call("service_hello", "hello.Hello", "SayHello", request, &output, context.Background()); err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("request:%v response:%v\n", request, output["message"])
		}
	}
}
