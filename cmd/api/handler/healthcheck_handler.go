package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterHealthCheck(mux *chi.Mux) {
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("\"pong\""))
	})
}
