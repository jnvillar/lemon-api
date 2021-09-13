package userservice

import (
	"context"

	usersmodel "lemonapp/domain/user/model"
	usersrepository "lemonapp/domain/user/repository"
)

type ServiceImpl struct {
	repository usersrepository.Repository
}

func NewServiceImpl(usersRepository usersrepository.Repository) Service {
	return &ServiceImpl{
		repository: usersRepository,
	}
}

func (s *ServiceImpl) List(ctx context.Context) ([]*usersmodel.User, error) {
	return s.repository.List(ctx)
}
func (s *ServiceImpl) Create(ctx context.Context, user *usersmodel.User) (*usersmodel.User, error) {
	return s.repository.Create(ctx, user)
}

func (s *ServiceImpl) GetByID(ctx context.Context, ID string) (*usersmodel.User, error) {
	return s.repository.GetByID(ctx, ID)
}
