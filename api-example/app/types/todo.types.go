package types

import "time"

type Todo struct {
	Name      string `json:"name"`
	IsChecked bool   `json:"isChecked"`
	Date      time.Time
}
