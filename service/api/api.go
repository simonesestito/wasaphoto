package api

import (
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/route"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Register a new route
	Register(route route.Route) error

	// Close terminates any resource used in the package
	Close() error
}

// NewRouter returns a new Router instance
func NewRouter(authMiddleware route.AuthMiddleware, middlewares []route.Middleware, logger logrus.FieldLogger) Router {
	// Create a new router where we will register HTTP endpoints.
	// The server will pass requests to this router to be handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{router, authMiddleware, middlewares, logger}
}

type _router struct {
	router         *httprouter.Router
	authMiddleware route.AuthMiddleware
	middlewares    []route.Middleware
	logger         logrus.FieldLogger
}

func (router _router) Handler() http.Handler {
	return router.router
}

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (router _router) Close() error {
	return nil
}

// Register a new route
func (router _router) Register(routeInfo route.Route) error {
	// Get route handler
	var handler route.Handler
	switch routeInfo.(type) {
	case route.AnonymousRoute:
		handler = routeInfo.(route.AnonymousRoute).Handler
	case route.SecureRoute:
		handler = router.authMiddleware.Intercept(routeInfo.(route.SecureRoute).Handler)
	default:
		return errors.New(fmt.Sprintf("Unknown route type: %s", reflect.TypeOf(routeInfo)))
	}

	// Wrap with other middlewares
	for _, middleware := range router.middlewares {
		handler = middleware(handler)
	}

	// Register path and method
	router.router.Handle(routeInfo.GetMethod(), routeInfo.GetPath(), router.wrap(handler))

	return nil
}

func (router _router) wrap(handle route.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			router.logger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ctx = route.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = router.logger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		handle(w, r, ps, ctx)
	}
}
