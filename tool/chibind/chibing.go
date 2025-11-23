package chibind

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
)

func DefaultBind[T any](r *http.Request, v *T) error {
	// Создание декодера
	decoder := form.NewDecoder()

	// ----- Parse request Body -----
	if err := json.NewDecoder(r.Body).Decode(v); err != nil && err != io.EOF {
		return err
	}

	// ----- Parse Url Params -----
	keys := chi.RouteContext(r.Context()).URLParams.Keys
	values := chi.RouteContext(r.Context()).URLParams.Values
	UrlParams := make(map[string][]string)
	for i := range keys {
		UrlParams[keys[i]] = []string{values[i]}
	}

	return decoder.Decode(v, UrlParams)
}
