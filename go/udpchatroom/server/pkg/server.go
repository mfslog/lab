package pkg

import (
	"github.com/mfslog/lab/go/udpchatroom/server/pkg/bucket"
	"github.com/mfslog/lab/go/udpchatroom/server/pkg/transport"
)

type server struct{
	bucket *bucket.Bucket
	transport *transport.Transport
}

var (
	DefaultServer = NewServer()
)


func NewServer()*server{
	return &server{
		bucket: bucket.NewBucket(),
	}
}

func (s server)Run(addr string)error{
	s.transport = transport.NewTransport(addr)
	return s.transport.Run()
}



