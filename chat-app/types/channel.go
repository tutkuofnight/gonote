package types

type Channel struct {
	Id       int       `gorm:"unique;primaryKey;autoIncrement"`
	Name     string    `json:"name" validate:"required,min=2"`
	AdminId  int       `json:"admin_id" gorm:"foreignKey:Id"`
	Messages []Message `json:"messages" gorm:"foreignKey:ChannelId"`
	Users    []*User   `json:"users" gorm:"many2many:user_channels;"`
}
