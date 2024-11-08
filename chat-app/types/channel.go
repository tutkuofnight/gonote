package types

type Channel struct {
	Id       int       `gorm:"unique;primaryKey;autoIncrement"`
	Name     string    `json:"name" validate:"required,min=2"`
	AuthorId int       `json:"authorId" gorm:"foreignKey:Id"`
	Messages []Message `json:"messages" gorm:"foreignKey:ChannelId"`
	Users    []User    `json:"users" gorm:"many2many:user_channels;"`
}
