package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Register a new route
	Register(route route.Route) error

	// RegisterAll controllers in the app, all at once, being library agnostic, as always
	RegisterAll(controllers []route.Controller) error

	// RegisterStatic handles requests for static files
	RegisterStatic(localPath string, httpPath string)

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

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (router *_router) Close() error {
	return nil
}
