package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imperfect-fourth/work"
	"github.com/imperfect-fourth/work/consumer"
	"github.com/imperfect-fourth/work/producer"
	"github.com/imperfect-fourth/work/transformer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func producerFn() ([]int, error) {
	fmt.Println("producing")
	return []int{0, 1, 2, 3, 4}, nil
}
func consumerFn(i int) error {
	fmt.Println(i)
	time.Sleep(1 * time.Second)
	return nil
}
func transformerFn(i int) (int, error) {
	time.Sleep(1 * time.Second)
	return i + 1, nil
}

func main() {
	traceProvider, err := startTracing()
	if err != nil {
		log.Fatalf("traceprovider: %v", err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("traceprovider: %v", err)
		}
	}()

	producer, producedIntChan, _ := work.NewProducer(
		"int producer",
		producerFn,
		producer.WithCooldown(10*time.Second),
	)
	go producer.Work()

	transformer, transformedIntChan, _ := work.NewTransformer(
		"int transformer",
		producedIntChan,
		transformerFn,
		transformer.WithWorkerPoolSize(1),
		transformer.WithSpanName("sleeping one and adding one"),
	)
	go transformer.Work()

	c, _ := work.NewConsumer("int consumer", transformedIntChan, consumerFn, consumer.WithSpanName("sleeping one and printing"))
	c.Work()

}

func startTracing() (*trace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	fmt.Println(os.Getenv("JAEGER_ENDPOINT"))
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(os.Getenv("JAEGER_ENDPOINT")),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("test"),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)

	return tracerprovider, nil
}
