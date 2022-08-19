package momolog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Log struct {
	stdout zerolog.Logger
	stderr zerolog.Logger
}

func New(options ...func(*config)) *Log {
	cfg := defaultConfig()
	for _, option := range options {
		option(&cfg)
	}

	stdout := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: cfg.formatTime}
	stderr := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: cfg.formatTime}

	return &Log{
		stdout: zerolog.New(stdout).Level(cfg.Level()).With().Timestamp().Logger(),
		stderr: zerolog.New(stderr).Level(cfg.Level()).With().Timestamp().Logger(),
	}
}

func (l *Log) Trace(ctx context.Context) *writer {
	return newWriter(ctx, l.stdout.Trace())
}

func (l *Log) Debug(ctx context.Context) *writer {
	return newWriter(ctx, l.stdout.Debug())
}

func (l *Log) Info(ctx context.Context) *writer {
	return newWriter(ctx, l.stdout.Info())
}

func (l *Log) Warn(ctx context.Context) *writer {
	return newWriter(ctx, l.stdout.Warn())
}

func (l *Log) Error(ctx context.Context) *writer {
	return newWriter(ctx, l.stderr.Error())
}

func (l *Log) Fatal(ctx context.Context) *writer {
	return newWriter(ctx, l.stderr.Fatal())
}
