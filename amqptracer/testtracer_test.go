package amqptracer

import (
	"strconv"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const testHTTPHeaderPrefix = "testprefix-"

// testTracer is a most-noop Tracer implementation that makes it possible for
// unittests to verify whether certain methods were / were not called.
type testTracer struct{}

var fakeIDSource = 1

func nextFakeID() int {
	fakeIDSource++
	return fakeIDSource
}

type testSpanContext struct {
	HasParent bool
	FakeID    int
}

func (n testSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

type testSpan struct {
	spanContext   testSpanContext
	OperationName string
	StartTime     time.Time
	Tags          map[string]interface{}
}

func (n testSpan) Equal(os opentracing.Span) bool {
	other, ok := os.(testSpan)
	if !ok {
		return false
	}
	if n.spanContext != other.spanContext {
		return false
	}
	if n.OperationName != other.OperationName {
		return false
	}
	if !n.StartTime.Equal(other.StartTime) {
		return false
	}
	if len(n.Tags) != len(other.Tags) {
		return false
	}

	for k, v := range n.Tags {
		if ov, ok := other.Tags[k]; !ok || ov != v {
			return false
		}
	}

	return true
}

// testSpan:
func (n testSpan) Context() opentracing.SpanContext                       { return n.spanContext }
func (n testSpan) SetTag(key string, value interface{}) opentracing.Span  { return n }
func (n testSpan) Finish()                                                {}
func (n testSpan) FinishWithOptions(opts opentracing.FinishOptions)       {}
func (n testSpan) LogFields(fields ...log.Field)                          {}
func (n testSpan) LogKV(kvs ...interface{})                               {}
func (n testSpan) SetOperationName(operationName string) opentracing.Span { return n }
func (n testSpan) Tracer() opentracing.Tracer                             { return testTracer{} }
func (n testSpan) SetBaggageItem(key, val string) opentracing.Span        { return n }
func (n testSpan) BaggageItem(key string) string                          { return "" }
func (n testSpan) LogEvent(event string)                                  {}
func (n testSpan) LogEventWithPayload(event string, payload interface{})  {}
func (n testSpan) Log(data opentracing.LogData)                           {}

// StartSpan belongs to the Tracer interface.
func (n testTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	sso := opentracing.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso)
	}
	return n.startSpanWithOptions(operationName, sso)
}

func (n testTracer) startSpanWithOptions(name string, opts opentracing.StartSpanOptions) opentracing.Span {
	fakeID := nextFakeID()
	if len(opts.References) > 0 {
		fakeID = opts.References[0].ReferencedContext.(testSpanContext).FakeID
	}

	return testSpan{
		OperationName: name,
		StartTime:     opts.StartTime,
		Tags:          opts.Tags,
		spanContext: testSpanContext{
			HasParent: len(opts.References) > 0,
			FakeID:    fakeID,
		},
	}
}

// Inject belongs to the Tracer interface.
func (n testTracer) Inject(sp opentracing.SpanContext, format interface{}, carrier interface{}) error {
	spanContext := sp.(testSpanContext)
	switch format {
	case opentracing.HTTPHeaders, opentracing.TextMap:
		carrier.(opentracing.TextMapWriter).Set(testHTTPHeaderPrefix+"fakeid", strconv.Itoa(spanContext.FakeID))
		return nil
	}
	return opentracing.ErrUnsupportedFormat
}

// Extract belongs to the Tracer interface.
func (n testTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	switch format {
	case opentracing.HTTPHeaders, opentracing.TextMap:
		// Just for testing purposes... generally not a worthwhile thing to
		// propagate.
		sm := testSpanContext{}
		err := carrier.(opentracing.TextMapReader).ForeachKey(func(key, val string) error {
			switch strings.ToLower(key) {
			case testHTTPHeaderPrefix + "fakeid":
				i, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				sm.FakeID = i
			}
			return nil
		})
		return sm, err
	}
	return nil, opentracing.ErrSpanContextNotFound
}
