package server

import (
	"github.com/mfslog/lab/go/udpchatroom/pkg/bucket"
	"github.com/mfslog/lab/go/udpchatroom/pkg/transport"
)

type server struct {
	bucket    *bucket.Bucket
	transport *transport.Transport
}

var (
	DefaultServer = NewServer()
)

func NewServer() *server {
	return &server{
		bucket: bucket.NewBucket(),
	}
}
