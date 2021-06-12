package logger

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/txfs19260817/scopelens/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateGlobalLogger builds a global zap logger
func CreateGlobalLogger() (logger *zap.Logger) {
	writeSyncer := getWriteSyncer()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Path.LogPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     60,
		Compress:   false,
	}
	return zapcore.AddSync(io.MultiWriter(os.Stdout, lumberJackLogger))
}
