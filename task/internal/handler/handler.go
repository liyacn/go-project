package handler

import (
	"project/pkg/logger"
	"project/pkg/storage"
	"project/pkg/wechat"
	"project/task/internal/service"
	"time"
)

type Config struct {
	Cos storage.CosConfig
	//Oss    storage.OssConfig
	Cdn    string
	Wechat struct {
		Appid  string
		Secret string
	}
	Robot struct {
		DingTalk   string
		WechatWork string
	}
}

type Handler struct {
	service     *service.Service
	storage     storage.API
	wechat      wechat.BasicAPI
	robotDing   string
	robotWechat string
}

func New(cfg *Config, srv *service.Service) *Handler {
	return &Handler{
		service: srv,
		storage: storage.NewCOS(&cfg.Cos),
		//storage:         storage.NewOSS(&cfg.Oss),
		wechat:      wechat.NewBasicAPI(cfg.Wechat.Appid, cfg.Wechat.Secret, logger.NewHttpClient(15*time.Second)),
		robotDing:   cfg.Robot.DingTalk,
		robotWechat: cfg.Robot.WechatWork,
	}
}
