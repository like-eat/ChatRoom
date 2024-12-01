package zlog

import (
	"apylee_chat_server/pkg/zlog"
	"go.uber.org/zap"
	"testing"
)

func TestInfo(t *testing.T) {
	zlog.Info("this is a info", zap.String("name", "apylee"))
}
