package main

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestSamplerFromEnvironmentHonorsParentBasedRatio(t *testing.T) {
	t.Setenv("OTEL_TRACES_SAMPLER", "parentbased_traceidratio")
	t.Setenv("OTEL_TRACES_SAMPLER_ARG", "0")

	sampler, err := samplerFromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	traceID := trace.TraceID{1}
	spanID := trace.SpanID{1}
	root := sampler.ShouldSample(sdktrace.SamplingParameters{TraceID: traceID})
	if root.Decision != sdktrace.Drop {
		t.Fatalf("expected an unsampled root at ratio zero, got %v", root.Decision)
	}

	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	child := sampler.ShouldSample(sdktrace.SamplingParameters{
		ParentContext: trace.ContextWithSpanContext(context.Background(), parent),
		TraceID:       traceID,
	})
	if child.Decision != sdktrace.RecordAndSample {
		t.Fatalf("expected a sampled parent to be inherited, got %v", child.Decision)
	}
}

func TestSamplerFromEnvironmentRejectsInvalidRatio(t *testing.T) {
	t.Setenv("OTEL_TRACES_SAMPLER", "parentbased_traceidratio")
	t.Setenv("OTEL_TRACES_SAMPLER_ARG", "1.5")

	if _, err := samplerFromEnvironment(); err == nil {
		t.Fatal("expected an invalid ratio to be rejected")
	}
}
