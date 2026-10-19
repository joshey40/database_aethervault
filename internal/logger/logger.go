package logger

import (
	"github.com/joshey40/database_aethervault/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger(config *config.Config) *zap.Logger {

	var lvl zap.AtomicLevel
	if config.Develop {
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	loggerCfg := zap.Config{
		Level:       lvl,
		Development: config.Develop,
		Sampling:    nil,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stdout",
		},
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			TimeKey:     "time",
			EncodeTime:  zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05"),
			LevelKey:    "lvl",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
		},
	}

	return zap.Must(loggerCfg.Build())

}
