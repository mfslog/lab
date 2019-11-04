package transport

import (
	"github.com/dubbogo/getty"
	"github.com/mfslog/lab/go/udpchatroom/pkg/pack"
	"time"
)

type netEventListener struct {
	Process pack.PackProcess
	session getty.Session
}

func NewEventLister() getty.EventListener {
	return &netEventListener{}
}

func (e *netEventListener) OnOpen(session getty.Session) error {
	e.session = session
	return nil
}

// invoked when session closed.
func (e *netEventListener) OnClose(session getty.Session) {

	return
}

// invoked when got error.
func (e *netEventListener) OnError(session getty.Session, err error) {

}

// invoked periodically, its period can be set by (Session)SetCronPeriod
func (e *netEventListener) OnCron(session getty.Session) {

}

//OnMessage cmd 消息响应
func (e *netEventListener) OnMessage(session getty.Session, udpCtx interface{}) {
	var (
		packPtr *pack.Pack
		ctx     getty.UDPContext
		ok      bool
	)

	ctx, ok = udpCtx.(getty.UDPContext)
	if !ok {
		return
	}

	packPtr, ok = ctx.Pkg.(*pack.Pack)
	if !ok {
		return
	}
	packPtr.PeerAddr = ctx.PeerAddr
	//
	if packPtr.Head.Response != nil ||
		packPtr.Head.Request != nil {
		e.Process(packPtr)
	}
}

func (e *netEventListener) WritePack(pack *pack.Pack) error {
	var (
		ctx getty.UDPContext
	)
	ctx.PeerAddr = pack.PeerAddr
	ctx.Pkg = pack

	return e.session.WritePkg(ctx, time.Second*5)
}
