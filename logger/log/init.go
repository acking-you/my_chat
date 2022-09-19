package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var logger = logrus.New()

//全局logger配置
func init() {
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			function = ""
			//处理文件名
			file = fmt.Sprintf(" %s:%d ", frame.File, frame.Line)
			return
		}},
	)
}

func Lg() *logrus.Logger {
	return logger
}
