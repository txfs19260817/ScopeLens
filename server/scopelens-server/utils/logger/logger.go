package logger

import (
	"scopelens-server/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

func init() {
	writeSyncer := getWriteSyncer()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	Logger, SugaredLogger = logger, logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Path.LogPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     60,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
