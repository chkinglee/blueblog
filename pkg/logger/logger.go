// Package logger
// @Author      : lilinzhen
// @Time        : 2022/2/20 20:12:48
// @Description :
package logger

import (
	"blueblog/internal/pkg/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultLevel      = zapcore.InfoLevel
	DefaultTimeLayout = time.RFC3339
)

type Option func(*option)

type option struct {
	level          zapcore.Level
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

func WithLevel(level string) Option {
	return func(o *option) {
		switch strings.ToLower(level) {
		case "debug":
			o.level = zapcore.DebugLevel
		case "info":
			o.level = zapcore.InfoLevel
		case "warn":
			o.level = zapcore.WarnLevel
		case "error":
			o.level = zapcore.ErrorLevel
		default:
			o.level = zapcore.InfoLevel
		}
	}
}

func WithFilePath(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	return func(o *option) {
		o.file = zapcore.Lock(f)
	}
}

func WithTimeLayout(timeLayout string) Option {
	return func(o *option) {
		o.timeLayout = timeLayout
	}
}

func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

func InfoT(ctx core.StdContext, msg string, fields ...zap.Field) {
	ctx.Logger.Info(msg, append(fields, zap.String("trace_id", ctx.Trace.ID()))...)
}

func DebugT(ctx core.Context, msg string, fields ...zap.Field) {
	ctx.Logger().Debug(msg, append(fields, zap.String("trace_id", ctx.Trace().ID()))...)
}

func WarnT(ctx core.Context, msg string, fields ...zap.Field) {
	ctx.Logger().Warn(msg, append(fields, zap.String("trace_id", ctx.Trace().ID()))...)
}

func ErrorT(ctx core.StdContext, msg string, fields ...zap.Field) {
	ctx.Logger.Error(msg, append(fields, zap.String("trace_id", ctx.Trace.ID()))...)
}

func NewJSONLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{level: DefaultLevel}
	for _, f := range opts {
		f(opt)
	}

	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	tee := zapcore.NewTee()

	if !opt.disableConsole {
		tee = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		tee = zapcore.NewTee(tee,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(tee,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	return logger, nil
}
