// Package tracing ...
package tracing

import (
	"gdemo/pcontext"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	SpanAttributeKeyTraceID    = "TraceID"
	SpanAttributeKeyController = "Controller"
	SpanAttributeKeyAction     = "Action"
	SpanAttributeKeyBody       = "Body"
	SpanAttributeKeyHeader     = "Header"

	SpanEventNameRequest  = "Request"
	SpanEventNameResponse = "Response"
)

type Span struct {
	trace.Span
}

func (s *Span) EndWithError(err error) {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
	}

	s.End()
}

func StartTrace(ctx *pcontext.Context, spanName string, opts ...trace.SpanStartOption) (*pcontext.Context, *Span) {
	c, span := otel.Tracer("").Start(ctx, spanName, opts...)
	span.SetAttributes(attribute.String(SpanAttributeKeyTraceID, ctx.TraceID()))

	return ctx.WithContext(c), &Span{span}
}

func StartTraceForFramework(ctx *pcontext.Context,
	spanName string, opts ...trace.SpanStartOption) (*pcontext.Context, trace.Span) {
	return StartTrace(ctx, spanName, opts...)
}
