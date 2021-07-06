package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"

	pkgHttp "go-workshop/pkg/http"
)


type Hello struct {
}

func (h *Hello) Register(ctx context.Context, router chi.Router) {
	router.Get("/hello", h.Hello)
}

func (h *Hello) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_ = pkgHttp.RequestJSONBody(w, r, http.StatusOK, map[string]interface{}{
		"Message": "Hello",
	})
}
