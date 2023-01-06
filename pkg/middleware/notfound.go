package middleware

import (
	"fmt"
	"net/http"

	"github.com/mozey/httprouter-util/pkg/response"
)

func NotFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(http.StatusNotFound, w, r, response.Response{
			Message: fmt.Sprintf("path not found %v", r.URL.Path),
		})
	})
}
