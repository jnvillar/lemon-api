package userservice

import (
	"context"

	usersmodel "lemonapp/domain/user/model"
)

type Service interface {
	Create(ctx context.Context, user *usersmodel.User) (*usersmodel.User, error)
	GetByID(ctx context.Context, ID string) (*usersmodel.User, error)
	List(ctx context.Context) ([]*usersmodel.User, error)
}
