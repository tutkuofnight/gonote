package repository

import (
	"api_example/app/types"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func FindByCredentials(users []types.User, username string, password string) (*types.User, error) {
	var findedUser *types.User
	for _, u := range users {
		if u.Username == username {
			findedUser = &u
		}
	}
	if findedUser != nil {
		status := bcrypt.CompareHashAndPassword([]byte(findedUser.Password), []byte(password))
		return findedUser, status
	} else {
		return findedUser, errors.New("user not found")
	}
}
