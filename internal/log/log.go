package log

import (
	"strings"

	"github.com/aredoff/reagate/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger       *zap.Logger
	globalLogger *zap.Logger
	zapConfig    zap.Config
)

func Init() {
	logLevel := config.Config.GetString("log.level")
	logFile := config.Config.GetString("log.file")

	if strings.ToLower(logLevel) == "debug" {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	if logFile != "" {
		zapConfig.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
		zapConfig.Encoding = "json"
		zapConfig.OutputPaths = []string{logFile}
		zapConfig.ErrorOutputPaths = []string{logFile}

	} else {
		zapConfig.OutputPaths = []string{"stderr"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.Encoding = "console"
	}

	Logger, _ = zapConfig.Build()
	zap.ReplaceGlobals(Logger)
	globalLogger = Logger.WithOptions(zap.AddCallerSkip(1))
}

func Warn(args ...interface{}) {
	globalLogger.Sugar().Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	globalLogger.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	globalLogger.Sugar().Errorf(template, args...)
}

func Error(args ...interface{}) {
	globalLogger.Sugar().Error(args...)
}

func Info(args ...interface{}) {
	globalLogger.Sugar().Info(args...)
}

func Infof(template string, args ...interface{}) {
	globalLogger.Sugar().Infof(template, args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Sugar().Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	globalLogger.Sugar().Fatalf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	globalLogger.Sugar().Debugf(template, args...)
}
