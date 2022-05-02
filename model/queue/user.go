package queue

/*
CustomTP nsq.topic
CustomCH nsq.channel
CustomEX rabbit.exchange
CustomQN rabbit.queue-name
Custom body-struct
*/

const (
	DefaultCH     = "-"
	AvatarToCdnTP = "avatar_to_cdn"
)

//const (
//	DefaultEX     = ""
//	AvatarToCdnQN = "avatar_to_cdn"
//)

type AvatarToCdn struct {
	ID        int    `json:"id"`
	Openid    string `json:"openid"`
	AvatarURL string `json:"avatar_url"`
}
