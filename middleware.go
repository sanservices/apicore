package apicore

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/san-services/apilogger"

	"github.com/dchest/uniuri"
)

const (
	APIKEY          string = "api-key"
	REQUEST_ID_KEY  string = "x-request-id"
	REMOTE_ADDR_KEY string = "remote-address"
	SESSIONID       string = "session"
)

var (
	errNoAPIKey     = errors.New("No api-key found")
	errAPIKeyInalid = errors.New("api-key not valid")
)

type responseWriterLogger struct {
	http.ResponseWriter
	statusCode int
}

func (rwl *responseWriterLogger) WriteHeader(code int) {
	rwl.statusCode = code
	rwl.ResponseWriter.WriteHeader(code)
}

// handleRequestInfo checks for requestId, api-key, clientIp and sessionId
// to add them to the context.
func handleRequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Add requestId to the request context.
		id := r.Header.Get(REQUEST_ID_KEY)

		if id == "" {
			id = uniuri.New()
		}

		rICtx := context.WithValue(r.Context(), REQUEST_ID_KEY, id)
		r = r.WithContext(rICtx)

		// Add client ip to the request context.
		addr := r.Header.Get("x-real-ip")

		if addr == "" {
			addr = r.RemoteAddr
		}

		adCtx := context.WithValue(r.Context(), REMOTE_ADDR_KEY, addr)
		r = r.WithContext(adCtx)

		// Add apiKey to the request context.
		apiKey := r.Header.Get(APIKEY)

		akCtx := context.WithValue(r.Context(), APIKEY, apiKey)
		r = r.WithContext(akCtx)

		// Add sessionId to the request context.
		sessionId := r.Header.Get(SESSIONID)

		sCtx := context.WithValue(r.Context(), SESSIONID, sessionId)
		r = r.WithContext(sCtx)

		next.ServeHTTP(w, r)
	})
}

// handleAPIKey is amiddleware to control the access to service
// in the presence of kong gateway don't use this.
func handleAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(string(APIKEY))
		lg := apilogger.New(r.Context(), "")

		if !strings.Contains(r.URL.Path, "docs") &&
			!strings.Contains(r.URL.Path, "health") && key == "" {

			lg.Error(apilogger.LogCatDebug, "No api-key in request")
			RespondError(r.Context(), w, http.StatusUnauthorized, errNoAPIKey)
			return
		}

		nCtx := context.WithValue(r.Context(), APIKEY, key)
		r = r.WithContext(nCtx)

		next.ServeHTTP(w, r)
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := apilogger.New(r.Context(), "")

		rwl := &responseWriterLogger{w, http.StatusOK}

		next.ServeHTTP(rwl, r)

		logFields := &apilogger.Fields{
			"status": rwl.statusCode,
			"url":    r.URL,
			"method": r.Method,
		}

		l.InfoWF(apilogger.LogCatReqPath, logFields)
	})
}
