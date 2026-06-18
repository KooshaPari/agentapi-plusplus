package observability

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestInitDisabled(t *testing.T) {
	t.Setenv(envExporter, "")

	shutdown, err := Init(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("init disabled observability: %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown disabled observability: %v", err)
	}
}

func TestInitStdout(t *testing.T) {
	t.Setenv(envExporter, "stdout")
	t.Setenv(envServiceName, "agentapi-test")

	shutdown, err := Init(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("init stdout observability: %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown stdout observability: %v", err)
	}
}

func TestInitRejectsUnsupportedExporter(t *testing.T) {
	t.Setenv(envExporter, "zipkin")

	if _, err := Init(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil))); err == nil {
		t.Fatal("expected unsupported exporter error")
	}
}
