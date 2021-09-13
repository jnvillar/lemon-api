package trasnfervalidator

import (
	"fmt"

	trasnfermodel "lemonapp/domain/transfer/model"
	"lemonapp/errors"
)

type Validator struct {
}

func NewTransferValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateTransfer(transfer *trasnfermodel.Transfer) error {
	if transfer == nil {
		return errors.NewBadRequestError(fmt.Errorf("must provide trasnfer"))
	}
	if transfer.Amount < 0 {
		return errors.NewBadRequestError(fmt.Errorf("transfer amount must be a positive int"))
	}
	if transfer.WalletTo == "" {
		return errors.NewBadRequestError(fmt.Errorf("transfer is missing wallet to"))
	}
	return nil
}

func (v *Validator) ValidateExtraction(transfer *trasnfermodel.Transfer) error {
	if transfer == nil {
		return errors.NewBadRequestError(fmt.Errorf("must provide trasnfer"))
	}
	if transfer.Amount < 0 {
		return errors.NewBadRequestError(fmt.Errorf("extraction amount must be a positive int"))
	}
	return nil
}

func (v *Validator) ValidateDeposit(transfer *trasnfermodel.Transfer) error {
	if transfer == nil {
		return errors.NewBadRequestError(fmt.Errorf("must provide trasnfer"))
	}
	if transfer.Amount < 0 {
		return errors.NewBadRequestError(fmt.Errorf("deposit amount must be a positive int"))
	}
	return nil
}
