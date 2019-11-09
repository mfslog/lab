//服务方法定义及分发

package server

import (
	"reflect"
	"sync"
)
type Handler interface {
	Name()string
	Handler() interface{}
}


type methodType struct{
	sync.Mutex
	method reflect.Method
	ArgType reflect.Type
	ReplyType reflect.Type
	ContextType reflect.Type
}


type service struct{
	name string
	rcvr reflect.Value
	typ  reflect.Type
	method map[string]*methodType
}


type dispatch struct{
	mu sync.Mutex
	serviceMap map[string]*service
}

func newDispatch()*dispatch{
	return &dispatch{
		mu:         sync.Mutex{},
		serviceMap: make(map[string]*service),
	}
}

func (d *dispatch)Handle(h Handler){

}

