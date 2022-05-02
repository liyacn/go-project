package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/api/internal/handler"
	"project/api/internal/service"
	"project/pkg/config"
	"project/pkg/logger"
	"project/pkg/process"
)

func initialize() *gin.Engine {
	path := flag.String("c", "conf.yaml", "config file")
	flag.Parse()
	cfg := config.Load[handler.Config, service.Config](*path)
	process.SetEnv(cfg.Env)
	logger.Initialize(&cfg.Logger)
	if process.GetEnv() != process.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}
	s := service.New(&cfg.Service)
	h := handler.NewEngine(&cfg.Handler, s)
	return h
}

func main() {
	h := initialize()
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
	l.Info("start")
	process.Notify()
	err := server.Shutdown(context.Background())
	process.DoCleanup()
	l.Info("stop", err)
}
