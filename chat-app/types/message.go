package types

import "time"

type Message struct {
	Text      string    `json:"text"`
	UserId    uint      `json:"user_id"`
	createdAt time.Time `gorm:"autoCreateTime"`
}
