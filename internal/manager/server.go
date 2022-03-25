package manager

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"jim_gateway/config"
	"net/http"
)

func RunPprof(config config.Config, ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", config.Pprof.Host,config.Pprof.Port)
	httpServer := http.Server{
		Addr:    address,
		Handler: http.DefaultServeMux,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				color.Red("close pprof")
				_ = httpServer.Shutdown(context.Background())
				return
			}
		}
	}()
	return httpServer.ListenAndServe()
}
