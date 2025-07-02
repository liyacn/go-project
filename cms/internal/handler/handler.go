package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/cms/internal/service"
	"project/model/cache"
	"project/pkg/ades"
	"project/pkg/captcha"
	"project/pkg/logger"
	"project/pkg/storage"
	"project/pkg/web/errcode"
	"slices"
	"strings"
)

type Config struct {
	Cos storage.CosConfig
	//Oss     storage.OssConfig
	Cdn     string
	Captcha string
	Aes     struct {
		Key string
		IV  string
	}
}

type Handler struct {
	service *service.Service
	storage storage.API
	cdn     string
	captcha string
	aes     ades.Cipher
	drawer  *captcha.Drawer
}

func New(cfg *Config, srv *service.Service) http.Handler {
	aes, err := ades.NewAesCTR([]byte(cfg.Aes.Key), []byte(cfg.Aes.IV))
	if err != nil {
		log.Fatal(err)
	}
	h := &Handler{
		service: srv,
		storage: storage.NewCOS(&cfg.Cos),
		//storage:     storage.NewOSS(&cfg.Oss),
		cdn:     cfg.Cdn,
		captcha: cfg.Captcha,
		aes:     aes,
		drawer:  captcha.NewDrawer("docs/fonts"),
	}
	r := gin.New()
	register(r, h)
	return r
}

const OK = http.StatusOK

func (h *Handler) CheckAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(errcode.Unauthorized.Response())
		return
	}
	auth, err := h.service.AdminTokenGet(c, token)
	if err != nil {
		logger.FromContext(c).Error("service.AdminTokenGet error", token, err)
		c.AbortWithStatusJSON(errcode.FromError(err))
		return
	}
	if auth.ID == 0 {
		c.AbortWithStatusJSON(errcode.Unauthorized.Response())
		return
	}
	c.Set("user", auth)
	c.Set("v3", auth.Username)
	c.Next()
}

func (h *Handler) RBAC(c *gin.Context) {
	auth := getAuth(c)
	if !auth.IsSuper && !slices.Contains(auth.Actions, c.Request.URL.Path) {
		c.AbortWithStatusJSON(errcode.PermissionDeny.Response())
		return
	}
	c.Next()
}

func (h *Handler) CdnLink(path string) string {
	if path == "" || strings.HasPrefix(path, "http") {
		return path
	}
	return h.cdn + path
}

func getAuth(c *gin.Context) *cache.AdminToken {
	v, _ := c.Get("user")
	return v.(*cache.AdminToken)
}
