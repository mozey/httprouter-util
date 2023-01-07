package middleware

import (
	"fmt"
	"net/http"

	"github.com/mozey/httprouter-util/pkg/handler"
	"github.com/mozey/httprouter-util/pkg/share"
)

func NotFound(h *handler.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.JSON(http.StatusNotFound, w, r, share.Response{
			Message: fmt.Sprintf("path not found %v", r.URL.Path),
		})
	})
}
