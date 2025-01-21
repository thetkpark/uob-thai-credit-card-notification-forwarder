package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME"`
	ENV         string `env:"ENV,default=local"`
	Level       string `env:"LOG_LEVEL,default=info"`
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}
func AppendCtxValue(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func Init(c Config, setAsDefault bool) *slog.Logger {
	defaultAttrs := []slog.Attr{
		slog.String("service", c.ServiceName),
		slog.String("env", c.ENV),
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     parseLevel(c.Level),
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts).WithAttrs(defaultAttrs)
	if c.ENV != "local" {
		handler = slog.NewJSONHandler(os.Stdout, opts).WithAttrs(defaultAttrs)
	}

	handler = ContextHandler{handler}

	logger := slog.New(handler)
	if setAsDefault {
		slog.SetDefault(logger)
	}
	return logger
}
