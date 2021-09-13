package userconverter

import (
	"encoding/json"
	"fmt"
	"net/http"

	usermodel "lemonapp/domain/user/model"
	customerrors "lemonapp/errors"
)

type Converter struct{}

func NewUsersConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertToUserFromRequest(request *http.Request) (*usermodel.User, error) {
	decoder := json.NewDecoder(request.Body)
	user := &usermodel.User{}
	err := decoder.Decode(user)
	if err != nil {
		return nil, customerrors.NewBadRequestError(fmt.Errorf("error parsing user %v", err))
	}
	return user, err
}
