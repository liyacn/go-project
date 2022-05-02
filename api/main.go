package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/api/internal/handler"
	"project/api/internal/service"
	"project/pkg/config"
	"project/pkg/logger"
	"project/pkg/process"
)

func setup() http.Handler {
	cfg := config.Load[handler.Config, service.Config]()
	process.SetEnv(cfg.App.Env)
	logger.Setup(&cfg.App.Logger)
	if process.GetEnv() != process.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}
	s := service.New(&cfg.Service)
	h := handler.New(&cfg.Handler, s)
	return h
}

func main() {
	h := setup()
	server := &http.Server{
		Addr:    ":8000",
		Handler: h,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	l := logger.New(process.GetHostname(), "api", process.GetIP(), process.GetEnv())
	l.Info("start", nil, nil)
	process.Notify()
	err := server.Shutdown(context.Background())
	process.DoCleanup()
	l.Info("stop", nil, err)
}
