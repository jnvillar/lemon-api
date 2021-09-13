// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"
	usermodel "lemonapp/domain/user/model"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *Service) Create(ctx context.Context, user *usermodel.User) (*usermodel.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *usermodel.User
	if rf, ok := ret.Get(0).(func(context.Context, *usermodel.User) *usermodel.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usermodel.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *usermodel.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, ID
func (_m *Service) GetByID(ctx context.Context, ID string) (*usermodel.User, error) {
	ret := _m.Called(ctx, ID)

	var r0 *usermodel.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *usermodel.User); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usermodel.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx
func (_m *Service) List(ctx context.Context) ([]*usermodel.User, error) {
	ret := _m.Called(ctx)

	var r0 []*usermodel.User
	if rf, ok := ret.Get(0).(func(context.Context) []*usermodel.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*usermodel.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
