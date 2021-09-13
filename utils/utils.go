package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	customerrors "lemonapp/errors"
	"lemonapp/logger"
)

func ReturnJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func ReturnError(w http.ResponseWriter, err error) {
	logger.Sugar.Errorf("error: %v", err)
	notFound := &customerrors.NotFoundError{}
	if errors.As(err, &notFound) {
		writeError(err.(*customerrors.NotFoundError).Code, w, err)
		return
	}
	badRequest := &customerrors.BadRequest{}
	if errors.As(err, &badRequest) {
		writeError(err.(*customerrors.BadRequest).Code, w, err)
		return
	}
	writeError(http.StatusInternalServerError, w, fmt.Errorf("internal server error"))
}

func writeError(status int, w http.ResponseWriter, err error) {
	w.WriteHeader(status)
	ReturnJSONResponse(w, &customerrors.Error{Error: err.Error()})
}
