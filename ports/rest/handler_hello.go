package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/valkyrie"
)

type FormBodyHello struct {
	Name  string `json:"name" validate:"required"`
	Greet string `json:"greet" validate:"required"`
}

type Hello struct {
}

func (h *Hello) Register(ctx context.Context, router *chi.Mux) {
	router.Post("/hello", bifrost.HandlerAdapter(h.Hello))
}

func (h *Hello) Hello(w http.ResponseWriter, r *http.Request) error {
	bifrost.JSONResponse(w)
	var form FormBodyHello
	if err := bifrost.RequestJSONBody(r, &form); err != nil {
		return bifrost.ErrMethodNotAllowed(w, r, err)
	}

	if err := valkyrie.Validate(form); err != nil {
		var errString []string
		for _, e := range err {
			errString = append(
				errString,
				fmt.Sprintf(
					"code: %s Type: %s Message: %s",
					e.Field, e.Type, e.Message,
				),
			)
		}
		return bifrost.ErrBadRequest(w, r, fmt.Errorf("%s", strings.Join(errString, ";\n")))
	}

	return bifrost.ResponsePayload(w, r, http.StatusOK, form)
}
