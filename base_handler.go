package apicore

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/san-services/apilogger"

	"github.com/gorilla/mux"
)

var (
	// ErrBadPayload the error to be returned when a request is malformed
	ErrBadPayload = errors.New("Malformed json request payload")
)

// BaseHandler is the handler for base api routes
type BaseHandler struct {
}

func newBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// InitRoutes initializes base (non-versioned) api routes
func (h *BaseHandler) InitRoutes(ctx context.Context, r *mux.Router) {
	r.HandleFunc("/healthcheck", h.healthCheck).Methods("GET")
}

func (h *BaseHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lg := apilogger.New(ctx, "")
	var out healthcheck
	var err error

	out.Host, err = os.Hostname()
	if err != nil {
		lg.Error(apilogger.LogCatHealth, err)
		RespondError(ctx, w, http.StatusOK, err)
		return
	}

	out.Datetime = time.Now()
	Respond(ctx, w, out, http.StatusOK, err)
}

// healthcheck response
type healthcheck struct {
	Host     string    `json:"host"`
	Datetime time.Time `json:"datetime"`
}
