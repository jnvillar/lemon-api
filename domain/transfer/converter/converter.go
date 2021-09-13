package transferconverter

import (
	"encoding/json"
	"fmt"
	"net/http"

	trasnfermodel "lemonapp/domain/transfer/model"
	customerrors "lemonapp/errors"
)

type Converter struct{}

func NewTransferConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertTransferFromRequest(request *http.Request) (*trasnfermodel.Transfer, error) {
	decoder := json.NewDecoder(request.Body)
	transfer := &trasnfermodel.Transfer{}
	err := decoder.Decode(transfer)
	if err != nil {
		return nil, customerrors.NewBadRequestError(fmt.Errorf("error parsing transfer %v", err))
	}
	return transfer, err
}
