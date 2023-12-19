// Package tracing ...
package tracing

import (
	"context"
	"crypto/rand"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"gdemo/pcontext"
)

type traceIdGenerator struct {
}

func (i *traceIdGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	rid := ctx.(*pcontext.Context).TraceID()

	uid := uuid.MustParse(rid)
	tid := trace.TraceID{}
	copy(tid[:], uid[:])

	return tid, i.NewSpanID(ctx, tid)
}

func (i *traceIdGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	sid := trace.SpanID{}
	_, _ = rand.Read(sid[:])

	return sid
}
