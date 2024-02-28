package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.NewAtomicLevel())
	logger = zap.New(core)
	defer logger.Sync()
}

func Info(msg string, fields ...string) {
	logger.Info(msg)
}

func Infof(msg string, format ...interface{}) {
	logger.Info(fmt.Sprintf(msg, format...))
}

func Errorf(msg string, format ...interface{}) {
	logger.Error(fmt.Sprintf(msg, format...))
}
