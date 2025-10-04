# Use the OpenTelemetry Collector Contrib as base image
FROM otel/opentelemetry-collector-contrib:latest

# Copy the collector configuration file
COPY configs/collector.yaml /etc/otelcol/config.yaml

# Override the default command to use our custom configuration
CMD ["--config=file:/etc/otelcol/config.yaml"]
