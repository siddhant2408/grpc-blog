package logging

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	return cfg.Build()
}
