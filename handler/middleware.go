package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterMiddlewares(r *mux.Router) {
	r.Use(addHeaders)
}

func addHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
