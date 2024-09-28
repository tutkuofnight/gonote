package types

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Todos    []Todo    `json:"todos"`
}
