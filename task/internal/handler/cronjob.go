package handler

import (
	"project/pkg/core"
	"project/pkg/logger"
	"time"
)

func (h *Handler) WechatTokenRefresh(ctx *core.Context) {
	l := logger.FromContext(ctx)
	ttl, err := h.service.WechatTokenTTL(ctx)
	if err != nil {
		l.Error("service.WechatTokenTTL error", err)
		return
	}
	if ttl > 10*time.Minute {
		return
	}
	resp, err := h.wechat.GetAccessToken(ctx)
	if err != nil {
		l.Error("wechat.AccessToken error", err)
		return
	}
	if resp.Errcode == 0 && resp.AccessToken != "" {
		err = h.service.WechatTokenSet(ctx, resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second)
		if err != nil {
			l.Error("service.WechatTokenSet error", err)
		}
	} else {
		l.Warn("wechat.AccessToken fail", resp)
	}
}
