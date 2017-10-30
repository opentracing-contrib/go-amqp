package amqptracer

import (
	"strconv"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

func TestAMQPHeaderInject(t *testing.T) {
	h := map[string]interface{}{}
	h["NotOT"] = "blah"
	h["opname"] = "AlsoNotOT"
	tracer := testTracer{}
	span := tracer.StartSpan("someSpan")
	fakeID := span.Context().(testSpanContext).FakeID

	// Use amqpHeadersCarrier to wrap around `h`.
	carrier := amqpHeadersCarrier(h)
	if err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier); err != nil {
		t.Fatal(err)
	}

	if len(h) != 3 {
		t.Errorf("Unexpected header length: %v", len(h))
	}
	// The prefix comes from just above; the suffix comes from
	// testTracer.Inject().
	if h["testprefix-fakeid"] != strconv.Itoa(fakeID) {
		t.Errorf("Could not find fakeid at expected key")
	}
}

func TestAMQPHeaderExtract(t *testing.T) {
	h := map[string]interface{}{}
	h["NotOT"] = "blah"
	h["opname"] = "AlsoNotOT"
	h["testprefix-fakeid"] = "42"
	tracer := testTracer{}

	// Use amqpHeadersCarrier to wrap around `h`.
	carrier := amqpHeadersCarrier(h)
	spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}

	if spanContext.(testSpanContext).FakeID != 42 {
		t.Errorf("Failed to read testprefix-fakeid correctly")
	}
}
