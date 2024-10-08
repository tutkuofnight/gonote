package types

type User struct {
	Id       string    `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Messages []Message `json:"messages"`
}
