package transport

import (
	"fmt"
	"github.com/dubbogo/getty"
	"github.com/mfslog/lab/go/intrepid/protocol/pack"
	"net"
	"time"
)

type tcpTransport struct{
	opts Options
}


func (t *tcpTransport)Read(session getty.Session, data []byte) (interface{}, int, error){
	intrepidPack := &pack.IntrepidPackage{}
	err := t.opts.Codec.Unmarshal(data,intrepidPack)
	return intrepidPack, len(data), err
}

func (t *tcpTransport)Write(session getty.Session, pack interface{}) ([]byte, error){
	data,err := t.opts.Codec.Marshal(pack)
	return data, err
}


type tcpTransportListener struct{
	l getty.Server
	tTransport *tcpTransport
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

		session.SetPkgHandler(t.tTransport)
		eventListener := NewTcpTransportSocket(session)
		session.SetEventListener(eventListener)
		fn(eventListener.(*tcpTransportSocket))
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


type tcpTransportClient struct{
	tt *tcpTransport
	addr string
	cc  getty.Client
	Socket
}







func (t *tcpTransport)Dial(addr string, opts ...DialOption)(Client, error){

	cc :=   getty.NewTCPClient(getty.WithServerAddress(addr))

	client := tcpTransportClient{
		tt:  t,
		addr: addr,
		cc:cc,
	}
	cc.RunEventLoop(func(session getty.Session) error {
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

		session.SetPkgHandler(t)
		eventListener := NewTcpTransportSocket(session)
		session.SetEventListener(eventListener)
		client.Socket = eventListener.(*tcpTransportSocket)
		return nil
	})

	return client, nil
}

func (t *tcpTransport)Listen(addr string, opts ...ListenOption)(Listener, error){
	var (
		options ListenOptions
	)

	for _, o := range opts{
		o(&options)
	}

	listener := &tcpTransportListener{
		l:          getty.NewTCPServer(getty.WithLocalAddress(addr)),
		tTransport: t,
	}
	return listener,nil
}