package apicore

import (
	"context"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/san-services/apilogger"
)

// initRouter initializes an API struct's router
func (api *API) initRouter(ctx context.Context) {
	lg := apilogger.New(ctx, "")
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.Use(handleRequestInfo)
	r.Use(requestLogger)

	for _, h := range api.handlers {
		h.InitRoutes(ctx, r)
	}

	api.Router = r

	if api.conf.Debug {
		err := api.Router.Walk(walkFn)
		if err != nil {
			lg.Error(apilogger.LogCatRouterInit, err)
		}
	}
}

// walkFn will print the initialized routes in a router
func walkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	if err != nil {
		return err
	}

	methods, _ := route.GetMethods()

	if len(methods) > 0 {
		pathInfo := fmt.Sprintf("%s %s", methods[0], path)
		log.Println(pathInfo)
	}
	return nil
}
