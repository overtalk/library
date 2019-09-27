package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is to called by other package
var Logger *zap.Logger

// todo： del later
func init() {
	InitLogger("info")
}

func InitLogger(lvString string) {
	var (
		err error
		lvl zapcore.Level
	)

	if lvl, err = getLoggerLevel(lvString); err != nil {
		log.Fatalln("failed to initialize logger due to:", err)
	}

	var loggerConfig zap.Config
	if lvl == zapcore.DebugLevel {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.Level = zap.NewAtomicLevelAt(lvl)
	Logger, err = loggerConfig.Build(
		zap.AddStacktrace(zapcore.PanicLevel),
	)

	if err != nil {
		log.Fatalln("failed to initialize logger due to:", err)
	}
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

// get log level
func getLoggerLevel(lvString string) (zapcore.Level, error) {
	var lvl zapcore.Level

	if err := lvl.UnmarshalText([]byte(lvString)); err != nil {
		return lvl, err
	}

	return lvl, nil
}
