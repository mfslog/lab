package pack

import (
	chat "github.com/mfslog/lab/go/udpchatroom/idl"
	"net"
)

type PackProcess func(pack *Pack)

type Pack struct {
	Head     *chat.Head
	PeerAddr *net.UDPAddr
	Body     []byte
}
