package rest

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/bifrost"
	"github.com/kubuskotak/valkyrie"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
)

type FormBodyHello struct {
	Name  string `json:"name" validate:"required"`
	Greet string `json:"greet" validate:"required"`
}

type Hello struct {
	Tracer opentracing.Tracer
}

func (h *Hello) Register(ctx context.Context, router chi.Router) {
	router.Use(bifrost.HttpTracer)
	router.Post("/hello", bifrost.HandlerAdapter(h.Hello))
}

func (h *Hello) Hello(w http.ResponseWriter, r *http.Request) error {
	span, _ := opentracing.StartSpanFromContext(r.Context(), "Hello.post")
	defer span.Finish()

	bifrost.JSONResponse(w)
	var form FormBodyHello
	if err := bifrost.RequestJSONBody(r, &form); err != nil {
		log.Error().Err(err).Msg("Request JSON body")
		return bifrost.ErrBadGateway(w, r, err)
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
		status := http.StatusBadRequest
		er := fmt.Errorf("%s", strings.Join(errString, ";\n"))
		ext.HTTPStatusCode.Set(span, uint16(status))
		ext.Error.Set(span, true)
		span.SetTag("error.kind", "validate")
		span.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
		span.SetTag("error.message", er.Error())
		span.SetTag("event", "hello.post.validate")
		span.LogKV(
			"message", er,
			"stack", string(debug.Stack()),
		)
		span.Finish()
		return bifrost.ErrBadRequest(w, r, er)
	}

	log.Info().Interface("Tracer ID", r.Context().Value(bifrost.TracerContext)).Msg("tracer id")

	return bifrost.ResponsePayload(w, r, http.StatusOK, form)
}
