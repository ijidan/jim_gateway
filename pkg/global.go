package pkg

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"io"
	"jim_message/config"
	"path/filepath"
	"runtime"
)

var (
	Root      string
	Conf      *config.Config
	Logger    *logrus.Logger
	Tracer    opentracing.Tracer
	Closer    io.Closer
	RequestId string
)

func Close() {

	_ = Closer.Close()
}
func init() {
	_, file, _, _ := runtime.Caller(0)
	Root = filepath.Dir(filepath.Dir(file))
	Conf = config.GetConfigInstance(Root)
	Logger = GetLoggerInstance(Conf, Root)
	Tracer, Closer = NewJaeger(*Conf, "jim_message")
	RequestId = "X-Request-Id"
}
