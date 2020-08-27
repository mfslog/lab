package envoy_tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type noopSpan struct{}
type noopSpanContext struct{}

var (
	defaultNoopSpanContext = noopSpanContext{}
	defaultNoopSpan        = noopSpan{}
	defaultNoopTracer      = EnvoyTracer{}
)

const (
	emptyString = ""
)

// noopSpanContext:
func (n noopSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

// noopSpan:
func (n noopSpan) Context() opentracing.SpanContext                                  { return defaultNoopSpanContext }
func (n noopSpan) SetBaggageItem(key, val string) opentracing.Span                   { return defaultNoopSpan }
func (n noopSpan) BaggageItem(key string) string                         { return emptyString }
func (n noopSpan) SetTag(key string, value interface{}) opentracing.Span             { return n }
func (n noopSpan) LogFields(fields ...log.Field)                         {}
func (n noopSpan) LogKV(keyVals ...interface{})                          {}
func (n noopSpan) Finish()                                               {}
func (n noopSpan) FinishWithOptions(opts opentracing.FinishOptions)                  {}
func (n noopSpan) SetOperationName(operationName string) opentracing.Span            { return n }
func (n noopSpan) Tracer() opentracing.Tracer                                        { return defaultNoopTracer }
func (n noopSpan) LogEvent(event string)                                 {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n noopSpan) Log(data opentracing.LogData)                                      {}
