package model

import (
	"golang.org/x/crypto/bcrypt"

	ierr "webChat/internal/errors"
)

type User struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Photo          []byte `json:"photo"`
	Email          string `json:"email"`
}

func (u *User) ComparePassword(password string) error {
	if bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)) != nil {
		return ierr.ErrIncorrectPassword
	}

	return nil
}
