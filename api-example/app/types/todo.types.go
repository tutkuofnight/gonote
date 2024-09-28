package types

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Name      string    `json:"name"`
	IsChecked bool      `json:"is_checked"`
	Date      time.Time
	//UserId    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
}
