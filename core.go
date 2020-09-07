package apicore

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/san-services/apicore/apisettings"
)

// API holds all
type API struct {
	conf     *apisettings.Service
	Router   *mux.Router
	handlers []Handler // v1.Handler, v2.Handler etc
}

// Handler is an interface that
// each api version's handler should implement
type Handler interface {
	InitRoutes(context.Context, *mux.Router)
}

// New constructs and returns a new API struct
func New(
	ctx context.Context,
	conf *apisettings.Service,
	handlers []Handler) *API {

	baseHandler := newBaseHandler()
	handlers = append(handlers, baseHandler)

	api := API{
		conf:     conf,
		handlers: handlers}

	api.initRouter(ctx)
	return &api
}
func ValidationRules() *validator.Validate {
	return validator.New()
}
