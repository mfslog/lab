package client

import "context"

type Client interface {
	Init(...Options)error
	Call(ctx context.Context, req, rsp interface{}, opts ...CallOptions)
	AsyncCall(ctx context.Context,req interface{}, fn func(rsp interface{}), opts ...CallOptions)
	Push(ctx context.Context, req interface{}, opts ...callOptions)
}
