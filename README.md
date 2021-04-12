# apicore

Package with set of helper functions to create a rest api.

# Installation

Run in yout terminal `go get -u github.com/sanservices/apicore`

# Description

There are different packages

## errs

* Package contains custom errors for api

examples:

```go
return helper.RespondError(c, http.StatusInternalServerError, errs.NewNoTemplateErr())
```

```go
return helper.RespondError(c, http.StatusInternalServerError, errs.NewExecTemplateErr())
```

## helper

Some basic set of helper functions

examples:

```go
return helper.RespondError(c, http.StatusBadRequest, errors.New("getExtendedInfo param must be boolean"))
```

```go
return helper.RespondOk(c, map[string]string{"ping":"pong"})
```

## middleware

Those middlewares should be added to echo router in every microservice

examples:

```go
package api

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	apicoreMW "github.com/sanservices/apicore/middleware"
	logger "github.com/sanservices/apilogger/v2"
)

// RegisterRoutes iterates over handlers and registers them in given echo server instance
func RegisterRoutes(e *echo.Echo, handlers []Handler) {
	e.Use(apicoreMW.SetCustomHeaders)
	e.Use(apicoreMW.EnrichContext)
	e.Use(apicoreMW.RequestLogger)
	e.Use(echoMW.Recover())

	for _, h := range handlers {
		h.RegisterRoutes(e.Group(""))
	}

	var routeList []string
	// Make an array of routes(formatted strings)
	for _, r := range e.Routes() {
		routeList = append(routeList, fmt.Sprintf(" [%s] %s ;", r.Method, r.Path))
	}
	logger.InfoWF(context.TODO(), logger.LogCatRouterInit, &logger.Fields{"routes": routeList})
}
```
