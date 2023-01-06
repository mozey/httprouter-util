package middleware

import (
	"fmt"
	"net/http"

	"github.com/mozey/httprouter-util/pkg/response"
	"github.com/pkg/errors"
)

type PanicHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

func PanicHandler() PanicHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		err := errors.WithStack(fmt.Errorf("%s", rcv))
		// response.JSON must print stack if resp type is error
		response.JSON(http.StatusInternalServerError, w, r, err)
	}
}
