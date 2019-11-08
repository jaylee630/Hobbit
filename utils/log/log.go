package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level       string `yaml:"level" json:"level"`
	LogFilePath string `yaml:"log_file_path" json:"log_file_path"`
}

type Logger struct {
	*zap.SugaredLogger
	cfg zap.Config
}

// Write make Logger can be used as a io.writer
func (l *Logger) Write(d []byte) (n int, err error) {

	l.Info(string(d))
	return len(d), nil
}

func NewLogger(cfg *LoggerConfig, opts ...zap.Option) (*Logger, error) {
	logger := &Logger{}
	logger.cfg = zap.NewProductionConfig()
	logger.cfg.Sampling = nil
	logger.cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// set log level
	if 0 != len(cfg.Level) {
		if err := logger.cfg.Level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return nil, fmt.Errorf("failed to parse logger level. error = %v, level = %s", err, cfg.Level)
		}
		if cfg.Level == "debug" {
			logger.cfg.Development = true
		}
	}

	// set log file path
	if 0 != len(cfg.LogFilePath) {
		logger.cfg.OutputPaths = []string{cfg.LogFilePath}
	}

	l, e := logger.cfg.Build()
	if e != nil {
		return nil, fmt.Errorf("failed to build logger. error = %v, cfg = %#v", e, logger.cfg)
	}

	if 0 < len(opts) {
		l = l.WithOptions(opts...)
	}

	logger.SugaredLogger = l.Sugar()
	return logger, nil
}

func (l *Logger) L() *zap.Logger {
	if nil == l.SugaredLogger {
		return zap.L()
	}
	return l.SugaredLogger.Desugar()
}

func (l *Logger) IsLogEnabled(level zapcore.Level) bool {
	return l.cfg.Level.Enabled(level)
}

func (l *Logger) SetLevel(s string) error {
	old := l.cfg.Level.String()
	if old != s {
		n := zap.NewAtomicLevel()
		if err := n.UnmarshalText([]byte(s)); err != nil {
			return fmt.Errorf("failed to change logger level. wrong level = %s", s)
		}

		if err := l.cfg.Level.UnmarshalText([]byte(s)); err != nil {
			return fmt.Errorf("failed to change logger level. wrong level = %s", s)
		}
		zap.S().Warnf("changed logger level from %s to %s.", old, s)
	}
	return nil
}
