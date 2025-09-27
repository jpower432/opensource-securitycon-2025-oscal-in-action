package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/complytime/complybeacon/proofwatch"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	olog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tracer      = otel.Tracer("conftest-exporter")
	meter       = otel.Meter("conftest-exporter")
	serviceName = semconv.ServiceNameKey.String("conftest-exporter")
)

func main() {
	var endpoint string
	var skipTLS, skipTLSVerify bool

	flag.StringVar(&endpoint, "endpoint", "", "The OTEL Collector Endpoint")
	flag.BoolVar(&skipTLS, "skip-tls", false, "Do not connect with TLS")
	flag.BoolVar(&skipTLSVerify, "skip-tls-verify", false, "Do not verify certificates")
	flag.Parse()

	conn, err := newClient(endpoint, true, true)
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err)
	}
	shutdown, err := otelSDKSetup(context.Background(), conn)
	if err != nil {
		log.Fatalf("failed to set up OpenTelemetry SDK: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down OpenTelemetry SDK: %v", err)
		}
	}()

	watcher, err := proofwatch.NewProofWatch("conftest-exporter", meter)
	if err != nil {
		log.Fatal(err.Error())
	}

	var results ConftestFileResult
	for _, finding := range results.Failures {
		evidence := finding.ToOCSF()
		if err := watcher.Log(context.Background(), evidence); err != nil {
			log.Fatalf("failed to log evidence %v", err)
		}
	}
}

// otelSDKSetup completes setup of the Otel SDK with providers.
func otelSDKSetup(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	shutDown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			serviceName,
		),
	)
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// And here, we set a global propagator. This is what handles injecting
	// context into gRPC metadata.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(3*time.Second))), sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	logProcessor := olog.NewSimpleProcessor(logExporter)
	logProvider := olog.NewLoggerProvider(olog.WithProcessor(logProcessor), olog.WithResource(res))

	// Register the provider as the global logger provider.
	global.SetLoggerProvider(logProvider)

	shutdownFuncs = append(shutdownFuncs, logProvider.Shutdown, meterProvider.Shutdown)

	return shutDown, nil
}

func newClient(otelEndpoint string, skipTLS, skipTLSVerify bool) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials
	if skipTLS {
		creds = insecure.NewCredentials()
	} else {
		sysPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to get system cert: %w", err)
		}
		// By default, skip TLS verify is false.
		creds = credentials.NewTLS(&tls.Config{RootCAs: sysPool, InsecureSkipVerify: skipTLSVerify}) /* #nosec G402  */ //pragma: allowlist secret
	}
	return grpc.NewClient(otelEndpoint, grpc.WithTransportCredentials(creds))
}
