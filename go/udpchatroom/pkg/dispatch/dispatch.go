package dispatch

import "github.com/mfslog/lab/go/udpchatroom/pkg/pack"

type Dispatch struct{
	clientMap map[string]pack.PackProcess
	SrvMap map[string]pack.PackProcess
}

var (
	DefaultDispatch = Dispatch{}
)

//func (d *Dispatch)