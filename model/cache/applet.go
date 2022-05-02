package cache

import (
	"strconv"
)

const (
	WechatTokenKey = "wx:tk" // 微信access_token
)

const (
	userTokenPre = "utk:"  // +token
	userInfoPre  = "user:" // +uid
)

func UserTokenKey(token string) string { return userTokenPre + token }

type UserToken struct {
	ID         int    `json:"i"`
	Openid     string `json:"o"`
	Unionid    string `json:"u"`
	SessionKey string `json:"s"`
}

func UserInfoKey(id int) string { return userInfoPre + strconv.Itoa(id) } // entity.User
