package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"project/api/internal/handler"
	"project/api/internal/service"
	"project/pkg/logger"
	"project/pkg/process"
)

func setup() (*service.Service, http.Handler) {
	var cfg struct {
		App struct {
			Env    string
			Logger logger.Config
		}
		Handler handler.Config
		Service service.Config
	}

	b, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}

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
	srv.Stop()
	l.Info("stop", nil, err)
}
