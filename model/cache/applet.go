package cache

import (
	"strconv"
)

const (
	WechatTokenKey = "wechat:token" // 微信access_token
)

const (
	userTokenPre = "user:token:" // +token
	userInfoPre  = "user:info:"  // +uid
)

func UserTokenKey(token string) string { return userTokenPre + token }

type UserToken struct {
	ID         int    `json:"i"`
	Openid     string `json:"o"`
	Unionid    string `json:"u"`
	SessionKey string `json:"s"`
}

func UserInfoKey(id int) string { return userInfoPre + strconv.Itoa(id) } // entity.User
