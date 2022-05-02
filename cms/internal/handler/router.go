package handler

import (
	"github.com/gin-gonic/gin"
	"project/pkg/web/middleware"
)

func register(r *gin.Engine, h *Handler) {
	r.GET("/", func(c *gin.Context) { c.String(OK, "OK") })
	r.Use(middleware.Cors, middleware.SetContext, middleware.Recover)

	r.POST("captcha", h.Captcha)
	{
		g := r.Group("", middleware.AccessLog)
		g.POST("login", h.Login)
		g.POST("logout", h.CheckAuth, h.Logout)
		g.POST("password", h.CheckAuth, h.Password)
	}

	rbacNoLog := r.Group("", h.CheckAuth, h.RBAC)
	rbac := rbacNoLog.Group("", middleware.AccessLog)

	{
		g := rbac.Group("system")
		g.POST("action/sync", h.SystemActionSync(r))
		g.POST("action/list", h.SystemActionList)
		g.POST("action/update", h.SystemActionUpdate)
		g.POST("role/list", h.AdminRoleList)
		g.POST("role/option", h.AdminRoleOption)
		g.POST("role/save", h.AdminRoleSave)
		g.POST("user/list", h.AdminUserList)
		g.POST("user/create", h.AdminUserCreate)
		g.POST("user/password", h.AdminUserPassword)
		g.POST("user/role", h.AdminUserRole)
		g.POST("user/status", h.AdminUserStatus)
	}

	{
		g := rbacNoLog.Group("upload")
		g.POST("image", h.UploadImage)
		g.POST("audio", h.UploadAudio)
		g.POST("video", h.UploadVideo)
		g.POST("pdf", h.UploadPDF)
	}

	{
		g := rbacNoLog.Group("crypto")
		g.POST("encrypt", h.CryptoEncrypt)
		g.POST("decrypt", h.CryptoDecrypt)
	}
}
