package grpcinvoke

import (
	"context"
	"errors"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

var (
	ErrProto3Invalid = errors.New("proto3 invalid")
)

//type FileMethodParam struct {
//	FileName   string           //文件名
//	MethodName string           //方法名
//	Params     *dynamic.Message //具体参数
//}
//
//type SetOpt func(o *Options)
//type Options struct {
//	kv      map[string]interface{}
//	context context.Context
//}
//
//func Context(ctx context.Context) SetOpt {
//	return func(o *Options) {
//		o.context = ctx
//	}
//}

/*
1、加载proto3文件（memory、file)
2、解析proto3格式 (json->inner []byte->inner)
3、填充proto3数据
4、序列化proto3格式化后的数据
5、invoke调用grpc服务返回结果
6、解析rpc返回的结果

*/

type ProtoSourceIf interface {
	load() error
	parse() error
	kv(filename, method, k string, v interface{})
	show()
	getservicehandler(filename, serivce, method string) *desc.MethodDescriptor
	getmessage(filename, method string) *dynamic.Message
}

type GrpcInvokeIf interface {
	loadProto(ProtoSourceIf) error
	parseProto() error
	fillData() error
	protoMarshal() ([]byte, error)
	call(filename, service, method string, req map[string]interface{}, rep *map[string]interface{}, ctx context.Context) error
}

type GrpcInvokeEngine struct {
	l ProtoSourceIf
	defaultInvoke
}

func (me *GrpcInvokeEngine) Parse() error {
	me.initGrpc()

	if err := me.loadProto(me.l); err != nil {
		return err
	}
	if err := me.parseProto(); err != nil {
		return err
	}

	if err := me.fillData(); err != nil {
		return err
	}

	return nil
}

func (me *GrpcInvokeEngine) Call(filename, servicename, methodname string,
	req map[string]interface{}, rep *map[string]interface{}, ctx context.Context) error {
	if err := me.call(filename, servicename, methodname, req, rep, ctx);
		err != nil {
		return err
	}
	return nil
}

func (me *GrpcInvokeEngine) SetKeyVal(filename, methodname, key string, value interface{}) {
	me.l.kv(filename, methodname, key, value)
}

func NewGrpcInvoke(loadder ProtoSourceIf) *GrpcInvokeEngine {
	return &GrpcInvokeEngine{
		l: loadder,
	}
}
