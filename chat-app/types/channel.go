package types

type Channel struct {
	Id       uint64    `json:"id"`
	Messages []Message `json:"messages"`
}
