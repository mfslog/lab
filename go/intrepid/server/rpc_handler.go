package server

import "reflect"

type rpcHandler struct{
	name string
	handler interface{}
}


func newRPCHandler(handler interface{})Handler{
	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	for m := 0; m < typ.NumMethod(); m++{

	}
}
