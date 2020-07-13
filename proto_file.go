package grpcinvoke

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"sync"
)

type protoMap struct {
	incl []string
	pm   map[string]string
	//proto3文件描述
	fd map[string]*desc.FileDescriptor
	//proto3[文件名]方法列表
	md  map[string]map[string]*dynamic.Message
	sds map[string][]string
	//屏蔽对外暴露
	rwmtx sync.RWMutex
}

type ProtoFile struct {
	protoMap
}

func (p *ProtoFile) ImportProtoWithList(filepaths map[string]string) {
	p.rwmtx.Lock()
	defer p.rwmtx.Unlock()
	p.pm = map[string]string{}
	for name, path := range filepaths {
		p.pm[name] = path
	}
}

func (p *ProtoFile) getservicehandler(filename, service, method string) *desc.MethodDescriptor {
	if fd, existfd := p.fd[filename]; existfd {
		for _, s := range fd.GetServices() {
			//fmt.Printf("servicename: %s pkgname:%s\n  input:%s", s.GetName(), s.GetFile().GetPackage(), service)
			if service == s.GetFile().GetPackage()+"."+s.GetName() {
				//fmt.Println("#########find md")
				return s.FindMethodByName(method)
			}
		}
	}

	return nil
}

func (p *ProtoFile) getmessage(filename, method string) *dynamic.Message {
	if fd, existfd := p.fd[filename]; existfd {
		for _, mds := range fd.GetMessageTypes() {
			//fmt.Printf("@@@ input method:%s mds:%s\n", method, mds.GetName())
			if mds.GetName() == method {
				return dynamic.NewMessage(mds)
			}
		}
	}
	return nil
}

func (p *ProtoFile) IncludeProto(path string) {
	if p.incl == nil {
		p.incl = []string{}
	}
	p.incl = append(p.incl, path)
}

func (p *ProtoFile) AddProto(name, path string) {
	if p.pm == nil {
		p.pm = make(map[string]string)
	}
	p.pm[name] = path
}


func (p *ProtoFile) load() error {
	p.rwmtx.RLock()
	defer p.rwmtx.RUnlock()
	p.fd = make(map[string]*desc.FileDescriptor)
	if p.pm == nil || len(p.pm) == 0 {
		return errors.New("proto files emtpy")
	}
	var ps protoparse.Parser
	ps.ImportPaths = p.incl
	//"/Users/luv/go/src/github.com/dbv/owl/protoset/test/test.proto"
	for name, path := range p.pm {
		myfds, err := ps.ParseFiles(path)
		if err != nil {
			return err
		}
		if len(myfds) != 1 {
			return fmt.Errorf("parse proto failed,len:", len(myfds))
		}
		p.fd[name] = myfds[0]
	}
	return nil
}

func (p *ProtoFile) kv(filename, method, key string, val interface{}) {
	if key == "" || val == nil {
		return
	}
	//fmt.Printf("set filename:%s method:%s key:%s val:%v\n", filename, method, key, val)
	if fd, fdexist := p.md[filename]; fdexist {
		if md, mdexist := fd[method]; mdexist {
			md.SetFieldByName(key, val)
			return
		}
	}
	return
}

func (p *ProtoFile) parse() error {
	p.md = make(map[string]map[string]*dynamic.Message)
	p.sds = make(map[string][]string)

	for filename, filedesc := range p.fd {
		p.md[filename] = make(map[string]*dynamic.Message)
		if !filedesc.IsProto3() {
			return ErrProto3Invalid
		}
		for _, service := range filedesc.GetServices() {
			p.sds[filename] = append(p.sds[filename], service.GetName())
		}
		methods := filedesc.GetMessageTypes()
		for _, method := range methods {
			p.md[filename][method.GetName()] = dynamic.NewMessage(method)
		}
	}

	//fmt.Println("service:", p.sds)
	//fmt.Println("md:", p.md)

	return nil
}

func (p *ProtoFile) show() {
	for filename, filedesc := range p.md {
		for methodname, methoddesc := range filedesc {
			jbuff, _ := methoddesc.Marshal()
			fmt.Printf("filename:%s method:%s desc:%v\n", filename, methodname, jbuff)

		}
	}
}
