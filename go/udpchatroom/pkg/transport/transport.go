package transport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dubbogo/getty"
	"github.com/gogo/protobuf/proto"
	chat "github.com/mfslog/lab/go/udpchatroom/idl"
	pack2 "github.com/mfslog/lab/go/udpchatroom/pkg/pack"
	"net"
	"time"
)

type Transport struct {
	l getty.Server
	eventListener getty.EventListener
}

func NewTransport(addr string)*Transport {
	return &Transport{
		l:             getty.NewUDPPEndPoint(getty.WithLocalAddress(addr)),
		eventListener: NewEventLister(),
	}
}

func (t *Transport)RegisterSrv(fn PackProcess){
	t.eventListener.(*netEventListener).Srv = fn
}

func (t *Transport)RegisterCli(fn PackProcess){
	t.eventListener.(*netEventListener).Cli = fn
}

func (t *Transport)Run()error{
	t.l.RunEventLoop(t.sessionCB)
	return nil
}

func (t *Transport)Read(session getty.Session,data  []byte)(interface{},int, error){
	buf := bytes.NewBuffer(data)
	totalLenBuf := buf.Next(4)
	totalLength := binary.BigEndian.Uint32(totalLenBuf)

	headLenBuf := buf.Next(4)
	headLength := binary.BigEndian.Uint32(headLenBuf)
	headData := buf.Next(int(headLength))

	bodyData := buf.Next(int(totalLength - headLength))
	pbHead := &chat.Head{}
	err := pbHead.Unmarshal(headData)
	if err != nil{
		return nil, 0, err
	}
	pack := &pack2.Pack{
		Head: pbHead,
		Body:bodyData,
	}

	return &pack, int(totalLength), nil
}

//Write 由于UDP 需要对端地址，所以这里
func (t *Transport)Write(session getty.Session, udpCtx interface{})([]byte,error){
	var (
		ctx *getty.UDPContext
		ok bool
		err error
		packPtr *pack2.Pack
		headData []byte
		bodyLen int
		headLen int
	)
	ctx, ok = udpCtx.(*getty.UDPContext)
	if !ok {
		return nil, errors.New("not UDPContext type")
	}
	//取包
	packPtr = ctx.Pkg.(*pack2.Pack)

	//序列化
	headData,err = proto.Marshal(packPtr.Head)

	bodyLen = len(packPtr.Body)
	headLen = len(headData)
	totalLenBuf := make([]byte,4,4)
	binary.BigEndian.PutUint32(totalLenBuf,uint32(headLen + bodyLen))

	headLenBuf := make([]byte, 4,4)
	binary.BigEndian.PutUint32(headLenBuf, uint32(headLen))

	buf := bytes.NewBuffer(make([]byte,headLen+bodyLen,headLen+bodyLen))

	buf.Write(totalLenBuf)
	buf.Write(headLenBuf)
	buf.Write(headData)
	buf.Write(packPtr.Body)

	return buf.Bytes(),err
}



func (t *Transport)sessionCB(session getty.Session)error{
	var (
		ok      bool
		udpConn *net.UDPConn
	)

	if udpConn, ok = session.Conn().(*net.UDPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not udp connection\n", session.Stat(), session.Conn()))
	}
	udpConn.SetReadBuffer(262144)
	udpConn.SetWriteBuffer(65536)

	session.SetName("udpServer")
	session.SetMaxMsgLen(1024)
	session.SetRQLen(1024)
	session.SetWQLen(1024)
	session.SetReadTimeout(time.Second)
	session.SetWriteTimeout(time.Second * 5)
	session.SetCronPeriod((int)(60 * time.Second))
	session.SetWaitTime(7 * time.Second)
	session.SetEventListener(t.eventListener)
	session.SetPkgHandler(t)
	return nil
}

