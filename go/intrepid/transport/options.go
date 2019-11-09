package transport

import (
	"context"
	"github.com/mfslog/lab/go/intrepid/protocol"
	"time"
)

type Options struct{
	Addrs []string
	Codec protocol.Codec
}


type DialOptions struct{
	Stream bool
	Timeout time.Duration
	Ctx context.Context
}

type DialOption func(o *DialOptions)


type ListenOptions struct{
	Ctx context.Context
}

type ListenOption func(o *ListenOptions)
