package handler

import (
	"bytes"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"io"
	"net/http"
	"project/model/queue"
	"project/pkg/core"
	"project/pkg/files"
	"project/pkg/logger"
)

func (h *Handler) AvatarToCdn(ctx *core.Context, msg *nsq.Message) error {
	var data queue.AvatarToCdn
	_ = json.Unmarshal(msg.Body, &data)
	ctx.Set("v2", data.Openid)
	l := logger.FromContext(ctx)
	resp, err := http.Get(data.AvatarURL)
	if err != nil {
		l.Error("http.Get error", data.AvatarURL, err)
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	ext := files.CheckImage(b)
	if ext == "" {
		l.Warn("invalid image type", data.AvatarURL)
		return nil
	}
	remotePath := "/avatar/" + files.GenFileName([]byte(data.Openid)) + "." + ext
	err = h.storage.PutObject(ctx, remotePath, bytes.NewReader(b))
	if err != nil {
		l.Error("storage.PutObject error", err)
		return err
	}
	err = h.service.UserAvatarUpdate(ctx, data.ID, remotePath)
	return err
}
