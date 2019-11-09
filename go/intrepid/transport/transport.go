package transport

import "github.com/mfslog/lab/go/intrepid/protocol/pack"

type Transport interface {
	Dial(addr string, opts ...DialOption)(Client, error)
	Listen(addr string, opts ...ListenOption)(Listener,error)
}

type Socket interface {
	Recv(intrepidPackage *pack.IntrepidPackage)
	Send(intrepidPackage *pack.IntrepidPackage)
	Close()error
	Local() string
	Remote() string
}

type Client interface {
	Socket
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(socket Socket))error
}
