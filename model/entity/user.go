package entity

const UserTable = "user"

type User struct {
	ID          int    `json:"-"`
	Openid      string `json:"-"`
	Unionid     string `json:"-"`
	PhoneNumber string `json:"phone_number"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
}

func (*User) TableName() string { return UserTable }
