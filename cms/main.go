package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/cms/internal/handler"
	"project/cms/internal/service"
	"project/pkg/config"
	"project/pkg/logger"
	"project/pkg/process"
)

func setup() (*service.Service, http.Handler) {
	cfg := config.Load[handler.Config, service.Config]()
	process.SetEnv(cfg.App.Env)
	logger.Setup(&cfg.App.Logger)
	if process.GetEnv() != process.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}
	s := service.New(&cfg.Service)
	h := handler.New(&cfg.Handler, s)
	return s, h
}

func main() {
	srv, h := setup()
	server := &http.Server{
		Addr:    ":6000",
		Handler: h,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	l := logger.New(process.GetHostname(), "cms", process.GetIP(), process.GetEnv())
	l.Info("start", nil, nil)
	process.Notify()
	err := server.Shutdown(context.Background())
	srv.Stop()
	l.Info("stop", nil, err)
}
