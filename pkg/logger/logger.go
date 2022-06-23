package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
)

const (
	DefaultLevel = logrus.DebugLevel
)

// NewLogger Create a Logrus logger with the nested-logrus-formatter format.
func NewLogger(field string) *logrus.Entry {
	exampleFormatter := &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return "", fmt.Sprintf("::%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			ShowFullLevel:   true,
			TrimMessages:    true,
			TimestampFormat: time.RFC822,
			HideKeys:        true,
		},
	}
	l := logrus.New()
	l.SetLevel(DefaultLevel)
	l.SetReportCaller(true)
	l.SetFormatter(exampleFormatter)

	return l.WithFields(logrus.Fields{"field": field})
}
