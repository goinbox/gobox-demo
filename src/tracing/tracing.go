// Package tracing ...
package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	traceresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace/noop"

	"gdemo/conf"
)

var tp *sdktrace.TracerProvider

func Init(config *conf.TracingConf, serviceName string) error {
	if !config.Enable {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil
	}

	res, err := traceresource.New(context.Background(),
		traceresource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceNamespaceKey.String(config.ServiceNamespace),
		),
	)
	if err != nil {
		return fmt.Errorf("traceresource.New error: %w", err)
	}

	exporter, err := genSpanExporter(config)
	if err != nil {
		return fmt.Errorf("genSpanExporter error: %w", err)
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(genSpanProcessor(config, exporter)),
		sdktrace.WithIDGenerator(&traceIdGenerator{}),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func Shutdown() error {
	if tp != nil {
		err := tp.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("tp.Shutdown error: %w", err)
		}
	}

	return nil
}

func genSpanExporter(config *conf.TracingConf) (sdktrace.SpanExporter, error) {
	switch config.Target {
	case "zipkin":
		return zipkin.New(config.ZipkinUrl)
	}

	return nil, fmt.Errorf("unsupport tracing target %s", config.Target)
}

func genSpanProcessor(config *conf.TracingConf, exporter sdktrace.SpanExporter) sdktrace.SpanProcessor {
	if config.Async {
		return sdktrace.NewBatchSpanProcessor(exporter)
	}
	return sdktrace.NewSimpleSpanProcessor(exporter)
}
