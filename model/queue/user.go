package queue

/*
CustomTP nsq.topic
CustomCH nsq.channel
Custom body-struct
*/

const (
	DefaultCH     = "-"
	AvatarToCdnTP = "avatar_to_cdn"
)

type AvatarToCdn struct {
	ID        int    `json:"id"`
	Openid    string `json:"openid"`
	AvatarURL string `json:"avatar_url"`
}
