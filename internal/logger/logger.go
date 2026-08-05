package logger

import "go.uber.org/zap"

var mainlogger *zap.Logger

func L() *zap.Logger {
	if mainlogger == nil {
		mainlogger, _ = zap.NewProduction()
	}
	return mainlogger
}
