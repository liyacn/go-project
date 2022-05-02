package entity

import "time"

type User struct {
	ID          int       `json:"id"`
	Openid      string    `json:"openid"`
	Unionid     string    `json:"unionid"`
	PhoneNumber string    `json:"phone_number"`
	Nickname    string    `json:"nickname"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at" gorm:"->"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"->"`
}

func (*User) TableName() string { return "user" }
