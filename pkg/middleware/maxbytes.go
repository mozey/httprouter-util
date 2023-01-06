package middleware

import "net/http"

type MaxBytesOptions struct {
	MaxBytes int64
}

// MaxBytes middleware can be used to limit POST body
// https://stackoverflow.com/a/28292505/639133
func MaxBytes(next http.Handler, o *MaxBytesOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, o.MaxBytes)
		next.ServeHTTP(w, r)
	})
}
