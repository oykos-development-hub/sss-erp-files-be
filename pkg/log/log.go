package log

/*
import (
	"log/slog"
	"sync"
)

var cfg config

type config struct {
	logger *slog.Logger
	sync.Mutex
}

/*
func Setup(app, env, version string) {
	var appAttr []slog.Attr
	if app != "" {
		appAttr = append(appAttr, App(app))
	}
	if env != "" {
		appAttr = append(appAttr, Env(env))
	}
	if version != "" {
		appAttr = append(appAttr, Version(version))
	}
	s := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			// AddSource: true,
		},
	).WithAttrs(appAttr)
	logHandler := NewContextHandler(s)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	cfg.logger = logger
}

//	LogFormat: &log.JSONFormatter{
//		FieldMap: log.FieldMap{
//			log.FieldKeyTime:  "@timestamp",
//			log.FieldKeyLevel: "log.level",
//			log.FieldKeyMsg:   "message",
//			log.FieldKeyFunc:  "function.name", // non-ECS
//		},
//	},
func Info(ctx context.Context, msg string, args ...any) {
	if cfg == (config{}) {
		return
	}
	cfg.logger.InfoContext(ctx, msg, args...)
}
func Warn(ctx context.Context, msg string, args ...any) {
	if cfg == (config{}) {
		return
	}
	cfg.logger.WarnContext(ctx, msg, args...)
}
func Error(ctx context.Context, msg string, args ...any) {
	if cfg == (config{}) {
		return
	}
	cfg.logger.ErrorContext(ctx, msg, args...)
}
*/
