package envoy_tracer

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type EnvoyTracer struct{}

var _ opentracing.Tracer = &EnvoyTracer{}

func (e EnvoyTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return defaultNoopSpan
}

func (e EnvoyTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return nil
}

func (e EnvoyTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	payload, ok := carrier.(map[string]string)
	logrus.Infof("judge result:%v", ok)
	if ok {
		for k, v := range payload {
			logrus.Info("key:%v value:%v", k, v)
		}
	}
	var niceMD metautils.NiceMD
	niceMD, ok = carrier.(metautils.NiceMD)
	if ok {
		val := niceMD.Get("x-b3-traceid")
		logrus.Info("x-b3-traceid", val)
	}

	var md opentracing.TextMapCarrier
	md, ok = carrier.(opentracing.TextMapCarrier)
	if ok {
		err := md.ForeachKey(func(key, val string) error {
			logrus.Infof("k:%v, val:%v", key, val)
			return nil
		})
		if err != nil {
			logrus.Errorf("err:%v", err)
		}
	}

	return nil, opentracing.ErrSpanContextNotFound
}
