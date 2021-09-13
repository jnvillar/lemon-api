package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Alias       string    `json:"alias"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"date_created"`
}

func NewUser(firstName, lastName, alias, email string) *User {
	return &User{
		ID:          uuid.NewString(),
		FirstName:   firstName,
		LastName:    lastName,
		Alias:       alias,
		Email:       email,
		DateCreated: time.Now().UTC(),
	}
}
