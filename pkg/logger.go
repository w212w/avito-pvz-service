package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger(level string) {
	Log.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		Log.Warn("Неизвестный уровень логирования, установлен уровень INFO")
		lvl = logrus.InfoLevel
	}
	Log.SetLevel(lvl)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.Debugf("Логгер инициализирован. Уровень: %s", lvl)
}
