package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {

	t.Run("simple logger", func(t *testing.T) {
		got := NewLogger("test")
		if got == nil {
			t.Error("nil logger")
			return
		}
		if got.Logger.Level != logrus.DebugLevel {

			t.Errorf("expected debug level, got %v", got.Level)
		}
	})

}
