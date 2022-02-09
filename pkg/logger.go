package pkg

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"io"
	"jim_gateway/config"
	"os"
	"sync"
)

var (
	onceLogger     sync.Once
	instanceLogger *logrus.Logger
)

func GetLoggerInstance(config *config.Config, root string) *logrus.Logger {
	onceLogger.Do(func() {
		logfile := root + "/" + config.Websocket.Log
		writer1 := colorable.NewColorableStdout()
		writer2, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		instanceLogger = logrus.New()
		instanceLogger.SetFormatter(&logrus.JSONFormatter{})
		instanceLogger.SetOutput(io.MultiWriter(writer1, writer2))
	})
	return instanceLogger
}
