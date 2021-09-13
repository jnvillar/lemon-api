package uservalidator

import (
	"fmt"
	"net/mail"

	usermodel "lemonapp/domain/user/model"
	"lemonapp/errors"
)

type Validator struct {
}

func NewUserValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateUser(user *usermodel.User) error {
	if user == nil {
		return errors.NewBadRequestError(fmt.Errorf("must provide user"))
	}
	if user.FirstName == "" {
		return errors.NewBadRequestError(fmt.Errorf("must provide firstname"))
	}
	if user.LastName == "" {
		return errors.NewBadRequestError(fmt.Errorf("must provide last name"))
	}
	if user.Alias == "" {
		return errors.NewBadRequestError(fmt.Errorf("must provide alias"))
	}
	if user.Email == "" {
		return errors.NewBadRequestError(fmt.Errorf("must provide email"))
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.NewBadRequestError(fmt.Errorf("invalid email address"))
	}
	return nil
}
