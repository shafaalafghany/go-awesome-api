package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level  string
	Output string
}

func NewZerolog(cfg Config) (zerolog.Logger, error) {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	logLevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Nop(), fmt.Errorf("invalid zerolog log level: %w", err)
	}
	zlog := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Logger()

	return zlog, nil
}
