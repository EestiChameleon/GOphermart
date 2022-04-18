package cmlogger

import (
	"go.uber.org/zap"
	"time"
)

var Sug *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("logger configured",
		// Structured context as loosely typed key-value pairs.
		"start time", time.Now().Format(time.RFC3339),
		"project", "GOphermart",
	)

	Sug = sugar
}
