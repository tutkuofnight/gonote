package types

type User struct {
	Username string    `json:"username"`
	Messages []Message `json:"messages"`
}
