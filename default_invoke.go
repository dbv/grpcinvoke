package grpcinvoke

import (
	"context"
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	methods map[string]dynamic.Message
}

type defaultInvoke struct {
	l  ProtoSourceIf
	cc *grpc.ClientConn //grpc链接

	md *desc.MessageDescriptor
	ms *dynamic.Message

	//req *InnerMessage
	//rep *InnerMessage
}

func (me *defaultInvoke) initGrpc() {
	var err error
	me.cc, err = grpc.Dial("172.30.30.46:31083", grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc connect failed:", err.Error())
		return
	}
}

func (me *defaultInvoke) loadProto(loadder ProtoSourceIf) error {
	me.l = loadder
	return me.l.load()
}

func (me *defaultInvoke) parseProto() error {
	//fmt.Println("start to pareProto")
	return me.l.parse()
}

func (me *defaultInvoke) fillData() error {
	//todo file check
	return nil
}

func (me *defaultInvoke) protoMarshal() ([]byte, error) {
	panic("implement me")
}

//service = package.servicename
func (me *defaultInvoke) call(filename, service, method string,
	req map[string]interface{}, rep *map[string]interface{}, ctx context.Context) error {
	//me.l.show()
	md := me.l.getservicehandler(filename, service, method)
	ms := me.l.getmessage(filename, md.GetInputType().GetName())

	if md == nil || ms == nil {
		fmt.Println("md or ms nil")
	}
	for key, value := range req {
		ms.SetFieldByName(key, value)
	}
	//ms.SetFieldByName("event", "myevent")
	//
	//fmt.Println("md.GetName():", md.GetName())
	//fmt.Println(ms.Marshal())

	if me.cc == nil {
		return errors.New("grpc conn nil")
	}
	stub := grpcdynamic.NewStub(me.cc)
	resp, err := stub.InvokeRpc(context.Background(), md, ms)
	if err != nil {
		return err
	}
	repMsg := resp.(*dynamic.Message)

	for key := range *rep {
		repfd := repMsg.GetMessageDescriptor().FindFieldByName(key)
		(*rep)[key] = repMsg.GetField(repfd)
	}
	return nil
}

func (me *defaultInvoke) loadProtoModel() error {
	return nil
}
