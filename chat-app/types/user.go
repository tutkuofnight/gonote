package types

//`gorm:"type:uuid;default:gen_random_uuid();primary_key;"`

type User struct {
	Id           int       `gorm:"unique;primaryKey;autoIncrement"`
	Username     string    `json:"username" validate:"required,min=2,max=20"`
	Password     string    `json:"password" validate:"required,min=6"`
	Channels     []Channel `json:"channels" gorm:"many2many:user_channels;"`
	ProfileImage string    `json:"profileImage" validate:"url"`
}
