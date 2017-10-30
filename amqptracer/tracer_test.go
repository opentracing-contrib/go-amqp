package amqptracer

import (
	"strconv"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

func TestInject(t *testing.T) {
	h := map[string]interface{}{}
	h["NotOT"] = "blah"
	h["opname"] = "AlsoNotOT"
	tracer := testTracer{}
	sp := tracer.StartSpan("someSpan")
	fakeID := sp.Context().(testSpanContext).FakeID

	// Inject the tracing context to the AMQP header.
	if err := Inject(sp, h); err != nil {
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

func TestExtract(t *testing.T) {
	h := map[string]interface{}{}
	h["NotOT"] = "blah"
	h["opname"] = "AlsoNotOT"
	h["testprefix-fakeid"] = "42"

	// Set the testTracer as the global tracer.
	opentracing.SetGlobalTracer(testTracer{})

	// Extract the tracing span out from the AMQP header.
	ctx, err := Extract(h)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.(testSpanContext).FakeID != 42 {
		t.Errorf("Failed to read testprefix-fakeid correctly")
	}
}
