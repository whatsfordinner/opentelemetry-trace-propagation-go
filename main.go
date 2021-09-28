package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Carrier implements the TextMapCarrier interface from the propagation
// module: https://pkg.go.dev/go.opentelemetry.io/otel/propagation#TextMapCarrier
// It can be marshalled into JSON for use in a polyglot environment
type Carrier struct {
	Fields map[string]string `json:"fields"`
}

func (c Carrier) Get(key string) string {
	return c.Fields[key]
}

func (c Carrier) Set(key string, value string) {
	c.Fields[key] = value
}

func (c Carrier) Keys() []string {
	keys := make([]string, 0, len(c.Fields))
	for key := range c.Fields {
		keys = append(keys, key)
	}

	return keys
}

// Message represents some sort of communication between processes unable
// to share a context. It might be placed onto a queue, emitted onto an
// event or do something truly mad scientist level
type Message struct {
	TraceContext Carrier `json:"traceContext"`
	Body         string  `json:"body"`
}

// initTracing sets up a new tracer provider using the stdout exporter
// Taken from: https://opentelemetry.io/docs/go/getting-started/
func initTracing() func(context.Context) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)

	if err != nil {
		panic("Unable to initialise stdout trace exporter")
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)

	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	otel.SetTextMapPropagator(propagator)

	return func(ctx context.Context) {
		_ = tracerProvider.Shutdown(ctx)
	}
}

// processA has context available to it but must pass work onto a separate
// process that it can only communicate with via a slice of bytes
func processA(ctx context.Context) {
	tracer := otel.Tracer("processA")
	ctx, span := tracer.Start(ctx, "processA")
	defer span.End()

	traceCarrier := Carrier{
		Fields: map[string]string{},
	}

	// Injecting the relevant key-value pairs for this trace into the
	// carrier. Because the carrier can be serialised into JSON it can be
	// passed to the other process without issues.
	otel.GetTextMapPropagator().Inject(ctx, traceCarrier)
	span.AddEvent(fmt.Sprintf("context injected into carrier: %+v", traceCarrier))

	newMessage := Message{
		TraceContext: traceCarrier,
		Body:         "a message meant for process B",
	}

	messageBytes, err := json.Marshal(newMessage)

	if err != nil {
		panic("Unable to marshal JSON message")
	}

	span.AddEvent("sending message to processB")
	processB(messageBytes)
	span.AddEvent("message sent to processB")
}

// processB doesn't receive the context from processA that contains the trace
// information - it must derive it from the contents of the byte slice that
// has been passed to it
func processB(message []byte) {
	tracer := otel.Tracer("processB")

	var newMessage Message
	err := json.Unmarshal(message, &newMessage)

	if err != nil {
		panic("Unable to unmarshal JSON message")
	}

	// Extracting the relevant values for continuing the trace from the carrier
	// and putting it into the local context.
	ctx := otel.GetTextMapPropagator().Extract(
		context.Background(),
		newMessage.TraceContext,
	)
	_, span := tracer.Start(ctx, "processB")
	span.AddEvent("context extracted from carrier, now we're in the right parent trace")
	defer span.End()
}

func main() {
	ctx := context.Background()
	shutdownTracing := initTracing()
	defer shutdownTracing(ctx)

	tracer := otel.Tracer("main")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	processA(ctx)
}
