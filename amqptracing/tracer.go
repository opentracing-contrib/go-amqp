package amqptracing

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

// Inject injects the span context into the AMQP header.
func Inject(span opentracing.Span, hdrs amqp.Table) error {
	c := amqpHeadersCarrier(hdrs)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, c)
}

// Extract extracts the span context out of the AMQP header.
func Extract(hdrs amqp.Table) (opentracing.SpanContext, error) {
	c := amqpHeadersCarrier(hdrs)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, c)
}
