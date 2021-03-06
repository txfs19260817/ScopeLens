package logger

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestCreateGlobalLogger(t *testing.T) {
	CreateGlobalLogger()
	zap.L().Error("computed sum = a+b",
		zap.Int("sum", 12),
		zap.Float32("a", 6.5),
		zap.Float64("b", 5.5),
		zap.String("url","https://www.elastic.co/guide/en/ecs-logging/go-zap/current/setup.html#setup-step-1"),
		zap.Error(fmt.Errorf("boom")))
}
