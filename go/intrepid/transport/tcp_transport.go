package transport

import (
	"fmt"
	"github.com/dubbogo/getty"
	"net"
	"time"
)

type tcpTransport struct{
	opts Options
}


type tcpTransportListener struct{
	l getty.Server
}

func(t *tcpTransportListener)Accept(fn func(socket Socket))error{
	t.l.RunEventLoop(func(session getty.Session) error {
		var (
			ok      bool
			tcpConn *net.TCPConn
		)

		if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
			panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
		}

		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(time.Duration(time.Second * 6))
		tcpConn.SetReadBuffer(262144)
		tcpConn.SetWriteBuffer(65536)

		session.SetName(fmt.Sprintf("tcp-%s", session.RemoteAddr()))
		session.SetMaxMsgLen(65536)
		session.SetRQLen(1024)
		session.SetWQLen(1024)
		session.SetReadTimeout(time.Second)
		session.SetWriteTimeout(time.Second * 5)
		session.SetCronPeriod(int(time.Second * 6))
		session.SetWaitTime(time.Second * 7)
		//session.SetTaskPool(t.taskPool)

		//session.SetPkgHandler(t)
		session.SetEventListener()
		return nil
	})
	return nil
}

func (t *tcpTransportListener)Addr()string{
	return t.l.Listener().Addr().String()
}

func (t *tcpTransportListener)Close()error{
	return t.l.Listener().Close()
}



func (t *tcpTransport)Dial(addr string, opts ...DialOption)(Client, error){

	return nil, nil
}

func (t *tcpTransport)Listen(addr string, opts ...ListenOption)(Listener, error){
	return nil,nil
}