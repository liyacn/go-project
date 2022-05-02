package handler

import (
	"github.com/gin-gonic/gin"
	"project/pkg/web/middleware"
)

func register(r *gin.Engine, h *Handler) {
	r.GET("/", func(c *gin.Context) { c.String(OK, "OK") })
	r.Use(middleware.SetContext, middleware.Recover)

	api := r.Group("", middleware.AccessLog)

	api.POST("wechat/login", h.WechatLogin)

	{
		g := api.Group("wechat", h.CheckAuth)
		g.POST("phone", h.WechatPhone)
		g.POST("user/profile", h.UserProfile)
		g.POST("user/info", h.UserInfo)
	}
}
