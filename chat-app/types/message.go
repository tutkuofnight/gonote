package types

import "time"

type Message struct {
	Id        int            `gorm:"primaryKey;autoIncrement"`
	Text      string         `json:"text" validate:"required"`
	ChannelId int            `json:"channel_id" gorm:"foreignKey:Id"`
	UserId    int            `json:"user_id" gorm:"foreignKey:Id"`
	createdAt time.Time      `gorm:"autoCreateTime"`
}

type MessageUserDto struct {
	Id           int `gorm:"foreignKey:Id"`
	Username     string `json:"username"`
	ProfileImage string `json:"profile_image"`
}

type MessageDto struct {
	Text string `json:"text"`
	User MessageUserDto
}
