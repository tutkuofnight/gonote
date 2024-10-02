package types

type User struct {
	Id           string `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Todos        []Todo `json:"todos" gorm:"foreignKey:UserId"`
}
