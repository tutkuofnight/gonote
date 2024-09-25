package types

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pin      int    `json:"pin"`
	Todos    []Todo `json:"todos"`
}
