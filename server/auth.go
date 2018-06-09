package server

import (
	"context"
	"net/http"
)

func authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User")

		ctx := context.WithValue(r.Context(), userContextKey, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
