package usermodel

import (
	walletmodel "lemonapp/domain/wallet/model"
)

type UserResponse struct {
	User    *User                `json:"user"`
	Wallets []walletmodel.Wallet `json:"wallets"`
}

func NewUserResponse(user *User, wallets []walletmodel.Wallet) *UserResponse {
	return &UserResponse{
		User: user, Wallets: wallets,
	}
}
