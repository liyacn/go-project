package queue

/*
CustomTP topic
CustomCH channel
CustomGP group
CustomEX exchange
CustomQN queue-name
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
