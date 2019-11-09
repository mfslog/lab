package transport

import (
	"github.com/dubbogo/getty"
	"github.com/mfslog/lab/go/intrepid/protocol/pack"
	"time"
)

type tcpTransportSocket struct{
	session getty.Session
	msgChan chan *pack.IntrepidPackage
}

func NewTcpTransportSocket(session getty.Session)getty.EventListener{
	return &tcpTransportSocket{
		session: session,
		msgChan: make(chan *pack.IntrepidPackage, 1024),
	}
}


func (t *tcpTransportSocket)Recv()*pack.IntrepidPackage{
	pkg := <-t.msgChan
	return pkg
}

func (t *tcpTransportSocket)Send(intrepidPackage *pack.IntrepidPackage){
	t.session.WritePkg(intrepidPackage, time.Second * 5)
}

func (t *tcpTransportSocket)Close()error{
	t.session.Close()
	return nil
}


func (t *tcpTransportSocket)Local() string{

	return t.session.LocalAddr()
}


func (t *tcpTransportSocket)Remote() string{
	return t.session.RemoteAddr()
}

func (t *tcpTransportSocket)OnOpen(session getty.Session)error{
	//t.session = session
	return nil
}

func (t *tcpTransportSocket)OnError(session getty.Session,err error){
	t.session.Close()
}


func (t *tcpTransportSocket)OnClose(session getty.Session){

}

func (t *tcpTransportSocket)OnMessage(session getty.Session, pkg interface{}){
	var (
		pbPkg  *pack.IntrepidPackage
		ok bool
	)
	pbPkg, ok = pkg.(*pack.IntrepidPackage)
	if !ok {
		return
	}
	t.msgChan<-pbPkg
}

func (t *tcpTransportSocket)OnCron(seession getty.Session){
}