package middleware

import (
	"context"
	"time"

	"github.com/dchest/uniuri"
	"github.com/labstack/echo/v4"
	logger "github.com/sanservices/apilogger/v2"
)

// SetCustomHeaders sets some custom headers
func SetCustomHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := next(c)

		id, ok := c.Request().Context().Value("x-request-id").(string)
		if !ok {
			id = ""
		}

		c.Response().Header().Set("X-Request-Id", id)

		return res
	}
}

// EnrichContext enriches the context and sets context to request
func EnrichContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()

		id := r.Header.Get(logger.RequestIDKey)
		key := r.Header.Get(logger.APIKEY)
		session := r.Header.Get(logger.SessionIDKey)

		if id == "" {
			id = uniuri.New()
		}

		addr := r.Header.Get("x-real-ip")
		if addr == "" {
			if r.RemoteAddr != "" {
				addr = r.RemoteAddr
			}
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, logger.APIKEY, key)
		ctx = context.WithValue(ctx, logger.RequestIDKey, id)
		ctx = context.WithValue(ctx, logger.RemoteAddrKey, addr)
		ctx = context.WithValue(ctx, logger.SessionIDKey, session)
		ctx = context.WithValue(ctx, logger.StartTime, time.Now())

		c.SetRequest(r.WithContext(ctx)) // set context to request and request to echo context

		return next(c)
	}
}

// RequestLogger logs requests
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()

		// Call next handler
		err := next(c)
		errCode := c.Response().Status
		httpErr, ok := err.(*echo.HTTPError)

		if ok {
			errCode = httpErr.Code
		}

		logger.InfoWF(r.Context(), logger.LogCatReqPath, &logger.Fields{
			"status": errCode,
			"url":    r.URL,
			"method": r.Method,
		})

		return err
	}
}
