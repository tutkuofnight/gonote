package note

import (
	"time"
)

type Note struct {
	Id       int
	Text     string
	Location map[string]interface{}
	Date     time.Time
}
