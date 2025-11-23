package logger

import (
	"fmt"
	"os"
	"storePrices/internal/platform/conf"
	"time"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ModuleName string

const (
	Product ModuleName = "product"
	Parser  ModuleName = "parser"
	Country ModuleName = "country"
	Worker  ModuleName = "worker"
	General ModuleName = "general"
)

type Factory interface {
	For(moduleName ModuleName) *zap.Logger
}

type factoryImpl struct {
	config    *conf.Config
	baseLevel zapcore.Level
}

func NewFactory(cfg *conf.Config) Factory {
	return &factoryImpl{
		config:    cfg,
		baseLevel: zap.InfoLevel,
	}
}

func (f *factoryImpl) For(moduleName ModuleName) *zap.Logger {
	moduleStr := string(moduleName)

	// Ensure the directory exists
	// If we don't do this, the app might crash on first run if the folder is missing.
	logDir := fmt.Sprintf("./logs/%s", moduleStr)
	_ = os.MkdirAll(logDir, 0755)

	// Dynamic Filename based on Date
	date := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%s/%s.log", logDir, date)

	// Configure Lumberjack (File rotation)
	fileRotator := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Define Writers
	var writers []zapcore.WriteSyncer

	// Write in files if tests no running
	if f.config.Env != "testing" {
		writers = []zapcore.WriteSyncer{zapcore.AddSync(fileRotator)}
	}

	// Show in console if NOT production
	if f.config.Env != "production" {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	coreWriter := zapcore.NewMultiWriteSyncer(writers...)

	// Encoder Config (JSON)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create Core
	core := zapcore.NewCore(encoder, coreWriter, f.baseLevel)

	// Build Logger
	return zap.New(core, zap.AddCaller()).With(zap.String("module", moduleStr))
}
