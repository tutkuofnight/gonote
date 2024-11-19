package types

type Channel struct {
	Id             string       `json:"id"`
	Name           string    `json:"name" validate:"required,min=2"`
	AuthorUsername string    `json:"authorUsername"`
	MaxMembers     int       `json:"maxMembers"`
	Messages       []Message `json:"messages" gorm:"foreignKey:ChannelId"`
	Users          []User    `json:"users" gorm:"many2many:user_channels;"`
}
