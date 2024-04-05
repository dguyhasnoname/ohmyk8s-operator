package util

import (
	"context"

	"go.uber.org/zap"

	"github.com/dguyhasnoname/ohmyk8s-operator/pkg/logs"
)

var (
	Logs *zap.SugaredLogger
	Ctx  context.Context
)

func init() {
	Ctx = context.Background()
	Logs = logs.Logger()
	Logs.Info("Initializing env...")
}
