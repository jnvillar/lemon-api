package transactionconverter

import (
	"fmt"
	"net/http"
	"strconv"

	currency "lemonapp/currency"
	transactionmodel "lemonapp/domain/transaction/model"
	"lemonapp/errors"
)

type Converter struct {
}

func NewTransactionConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertSearchParamsFromRequest(request *http.Request) (*transactionmodel.SearchParams, error) {
	params := request.URL.Query()
	var err error
	var limit, offset int
	var transactionType transactionmodel.TransactionType
	var curr currency.Currency

	if len(params["limit"]) > 0 {
		limit, err = strconv.Atoi(params["limit"][0])
		if err != nil {
			return nil, errors.NewBadRequestError(fmt.Errorf("invalid limit"))
		}
	}

	if len(params["offset"]) > 0 {
		offset, err = strconv.Atoi(params["offset"][0])
		if err != nil {
			return nil, errors.NewBadRequestError(fmt.Errorf("invalid offset"))
		}
	}

	if len(params["transaction_type"]) > 0 {
		transactionType = params["transaction_type"][0]
		if valid := transactionmodel.ValidTransaction(transactionType); !valid {
			return nil, errors.NewBadRequestError(fmt.Errorf("unknown transaction type"))
		}
	}

	if len(params["currency"]) > 0 {
		curr = params["currency"][0]
		if valid := currency.ValidCurrency(curr); !valid {
			return nil, errors.NewBadRequestError(fmt.Errorf("unknown currency"))
		}
	}

	return &transactionmodel.SearchParams{
		Limit:           limit,
		Offset:          offset,
		TransactionType: transactionType,
		Currency:        curr,
	}, nil
}
