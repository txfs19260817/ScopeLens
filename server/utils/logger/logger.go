package logger

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/txfs19260817/scopelens/server/config"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateGlobalLogger builds a global zap logger
func CreateGlobalLogger() (logger *zap.Logger) {
	logger = zap.New(ecszap.WrapCore(zapcore.NewCore(getEncoder(), getWriteSyncer(), zapcore.DebugLevel)), zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := ecszap.EncoderConfig{
		EnableName:       true,
		EnableCaller:     true,
		EnableStackTrace: true,
		EncodeName:       zapcore.FullNameEncoder,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     ecszap.ShortCallerEncoder,
	}.ToZapCoreEncoderConfig()
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
