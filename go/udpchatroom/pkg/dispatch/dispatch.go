package dispatch

import (
	"context"
	"github.com/mfslog/lab/go/udpchatroom/pkg/pack"
	"time"
)

type MethodFunc func(ctx context.Context, req []byte) ([]byte, error)
type Dispatch struct {
	clientMap map[string]MethodFunc
	SrvMap    map[string]MethodFunc
}

var (
	DefaultDispatch = &Dispatch{
		clientMap: make(map[string]MethodFunc),
		SrvMap:    make(map[string]MethodFunc),
	}
)

func (d *Dispatch) ProcessPack(pack *pack.Pack) {
	var (
		fn     MethodFunc
		ok     bool
		ctx    context.Context
		cancel context.CancelFunc
		resp   []byte
		err    error
	)

	if pack.Head.Request != nil {
		if fn, ok = d.SrvMap[pack.Head.Request.MethodName]; !ok {
			return
		}
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	go func() {
		resp, err = fn(ctx, pack.Body)
	}()

	select {
	case <-ctx.Done():

	case <-time.After(time.Second * 5):

	}
}
