package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

type QuestifyLogger struct {
	logger *slog.Logger
}

func New(debug bool) *QuestifyLogger {
	lv := slog.LevelInfo
	if debug {
		lv = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: lv}
	hdr := slog.NewJSONHandler(os.Stdout, opts)
	lgr := slog.New(hdr)

	return &ZionLogger{lgr}
}

func (l *ZionLogger) Fatal(msg string) {
	l.Logger.Error(msg)
	os.Exit(1)
}
