package types

import "time"

type Message struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	Text      string    `json:"text" validate:"required"`
	ChannelId int       `json:"channelId" gorm:"foreignKey:Id"`
	UserId    int       `json:"userId" gorm:"foreignKey:Id"`
	createdAt time.Time `gorm:"autoCreateTime"`
}

type MessageUserDto struct {
	Id           int    `gorm:"foreignKey:Id"`
	Username     string `json:"username"`
	ProfileImage string `json:"profileImage"`
}

type MessageDto struct {
	Id   int            `json:"id"`
	Text string         `json:"text"`
	User MessageUserDto `json:"user"`
}
