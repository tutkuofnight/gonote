package types

import "time"

type Message struct {
	Id        int       `gorm:"unique;primaryKey;autoIncrement"`
	Text      string    `json:"text" validate:"required"`
	ChannelId int       `json:"channel_id" gorm:"foreignKey:Id"`
	UserId    int       `json:"user_id" gorm:"foreignKey:Id"`
	createdAt time.Time `gorm:"autoCreateTime"`
}
