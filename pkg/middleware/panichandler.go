package middleware

import (
	"fmt"
	"net/http"

	"github.com/mozey/httprouter-util/pkg/handler"
	"github.com/pkg/errors"
)

type PanicHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

func PanicHandler(h *handler.Handler) PanicHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		err := errors.WithStack(fmt.Errorf("%s", rcv))
		// response.JSON must print stack if resp type is error
		h.JSON(http.StatusInternalServerError, w, r, err)
	}
}
