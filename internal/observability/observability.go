package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	envExporter     = "AGENTAPI_OTEL_EXPORTER"
	envServiceName  = "OTEL_SERVICE_NAME"
	envOTLPEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"
)

// Shutdown releases OpenTelemetry SDK resources.
type Shutdown func(context.Context) error

// Init installs OpenTelemetry providers when AGENTAPI_OTEL_EXPORTER is set to
// "otlp" or "stdout". The default is "none" to preserve existing CLI output.
func Init(ctx context.Context, logger *slog.Logger) (Shutdown, error) {
	exporterName := strings.ToLower(strings.TrimSpace(os.Getenv(envExporter)))
	if exporterName == "" || exporterName == "none" {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create otel resource: %w", err)
	}

	traceExporter, err := newTraceExporter(ctx, exporterName)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("Initialized OpenTelemetry", "exporter", exporterName, "service", serviceName())

	return func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var shutdownErr error
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			shutdownErr = fmt.Errorf("shutdown tracer provider: %w", err)
		}
		return shutdownErr
	}, nil
}

func newTraceExporter(ctx context.Context, exporterName string) (sdktrace.SpanExporter, error) {
	switch exporterName {
	case "otlp":
		return newOTLPHTTPExporter(), nil
	case "stdout":
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("create stdout trace exporter: %w", err)
		}
		return exporter, nil
	default:
		return nil, fmt.Errorf("unsupported %s value %q (want otlp, stdout, or none)", envExporter, exporterName)
	}
}

type otlpHTTPExporter struct {
	client   *http.Client
	endpoint string
}

func newOTLPHTTPExporter() *otlpHTTPExporter {
	endpoint := strings.TrimRight(os.Getenv(envOTLPEndpoint), "/")
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}
	return &otlpHTTPExporter{
		client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: endpoint + "/v1/traces",
	}
}

func (e *otlpHTTPExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	if len(spans) == 0 {
		return nil
	}

	payload := map[string]any{
		"resourceSpans": []map[string]any{{
			"scopeSpans": []map[string]any{{
				"spans": spansForOTLP(spans),
			}},
		}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal OTLP spans: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create OTLP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("export OTLP spans: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("export OTLP spans: status %s", resp.Status)
	}
	return nil
}

func (e *otlpHTTPExporter) Shutdown(context.Context) error {
	return nil
}

func spansForOTLP(spans []sdktrace.ReadOnlySpan) []map[string]any {
	out := make([]map[string]any, 0, len(spans))
	for _, span := range spans {
		spanContext := span.SpanContext()
		out = append(out, map[string]any{
			"traceId":           spanContext.TraceID().String(),
			"spanId":            spanContext.SpanID().String(),
			"parentSpanId":      span.Parent().SpanID().String(),
			"name":              span.Name(),
			"kind":              int(span.SpanKind()),
			"startTimeUnixNano": span.StartTime().UnixNano(),
			"endTimeUnixNano":   span.EndTime().UnixNano(),
		})
	}
	return out
}

func serviceName() string {
	if name := strings.TrimSpace(os.Getenv(envServiceName)); name != "" {
		return name
	}
	return "agentapi"
}
