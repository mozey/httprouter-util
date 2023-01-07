package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/alecthomas/units"
	"github.com/pkg/errors"
)

// GetBody from http.Request
func (h *Handler) GetBody(r *http.Request) (body []byte, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		body = []byte("")
	} else {
		maxPayload, err := h.Config.FnMaxPayloadMb().Int64()
		if err != nil {
			return body, errors.WithStack(err)
		}
		body, err = ioutil.ReadAll(
			io.LimitReader(r.Body, maxPayload*int64(units.MiB)))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := r.Body.Close(); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return body, nil
}
