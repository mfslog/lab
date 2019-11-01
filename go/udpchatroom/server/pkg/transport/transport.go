package transport

import (
	"bytes"
	"encoding/binary"
	"github.com/dubbogo/getty"
	chat "github.com/mfslog/lab/go/udpchatroom/idl"
	"github.com/sirupsen/logrus"
)

type Transport struct {
	l getty.Server
}


type Pack struct{
	Meta *chat.Meta
	Data []byte
}



func NewTransport(addr string)*Transport {
	return &Transport{}
}

func (t *Transport)Run()error{
	t.l.RunEventLoop(t.sessionCB)
	return nil
}

func (t *Transport)Read(session getty.Session,data  []byte)(interface{},int, error){
	buf := bytes.NewBuffer(data)
	totalLenBuf := buf.Next(4)
	totalLength := binary.BigEndian.Uint32(totalLenBuf)

	metaLenBuf := buf.Next(4)
	metaLen := binary.BigEndian.Uint32(metaLenBuf)
	metaData := buf.Next(int(metaLen))

	msgData := buf.Next(int(totalLength - metaLen))
	meta := &chat.Meta{}
	err := meta.Unmarshal(metaData)
	if err != nil{
		return nil, 0, err
	}
	pack := &Pack{
		Meta: meta,
		Data:msgData,
	}

	return &pack, int(totalLength), nil
}

func (t *Transport)Write(session getty.Session, udpCtx interface{})([]byte,error){
	ctx := udpCtx.(getty.UDPContext)
	pack := ctx.Pkg.(Pack)
	metaData, err :=pack.Meta.Marshal()
	if err != nil{
		return nil,err
	}
	//todo: to pool
	metaLen := len(metaData)
	DataLen := len(pack.Data)
	totalLen := metaLen + DataLen
	totalLenBuf := make([]byte,4,4)
	binary.BigEndian.PutUint32(totalLenBuf,uint32(totalLen))

	metaLenBuf := make([]byte, 4,4)
	binary.BigEndian.PutUint32(metaLenBuf, uint32(metaLen))

	buf := bytes.NewBuffer(make([]byte,totalLen,totalLen))

	buf.Write(metaLenBuf)
	buf.Write(metaData)
	buf.Write(pack.Data)
	return buf.Bytes(),nil
}



func (t *Transport)sessionCB(session getty.Session)error{
	session.SetEventListener(t)
	session.SetPkgHandler(t)
	return nil
}

//OnOpen UDP 会话只有一个session
func (t *Transport) OnOpen(session getty.Session) error {
	logrus.Infof("got session:%s", session.Stat())
	return nil
}

func (t *Transport) OnError(session getty.Session, err error) {
	logrus.Infof("session{%s} got error{%v}, will be closed.", session.Stat(), err)

}

func (t *Transport) OnClose(session getty.Session) {
	logrus.Infof("session{%s} is closing......", session.Stat())

}

func (t *Transport) OnMessage(session getty.Session, udpCtx interface{}) {


}

func (e *Transport) OnCron(session getty.Session) {

}

