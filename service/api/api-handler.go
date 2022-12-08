package api

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

func (router *_router) RegisterAll(controllers []route.Controller) error {
	// Register routes
	for _, controller := range controllers {
		for _, routeInfo := range controller.ListRoutes() {
			if err := router.Register(routeInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

func (router *_router) Handler() http.Handler {
	return router.router
}

// RegisterStatic handles requests for static files
func (router *_router) RegisterStatic(localPath string, httpPath string) {
	router.logger.Debugln("Registering static path", localPath, "to HTTP path", httpPath)
	router.router.ServeFiles(httpPath+"/*filepath", http.Dir(localPath))
}

// Register a new route
func (router *_router) Register(routeInfo route.Route) error {
	// Get route handler
	var handler route.Handler

	anonymousRoute, isAnonymous := routeInfo.(route.AnonymousRoute)
	secureRoute, isSecure := routeInfo.(route.SecureRoute)

	switch {
	case isAnonymous:
		handler = anonymousRoute.Handler
	case isSecure:
		handler = router.authMiddleware.Intercept(secureRoute.Handler)
	default:
		return fmt.Errorf("Unknown route type: %s", reflect.TypeOf(routeInfo))
	}

	// Wrap with other middlewares
	for _, middleware := range router.middlewares {
		handler = middleware(handler)
	}

	// Register path and method
	router.logger.Debugf(
		"Registering route of type '%s' [%s] %s",
		reflect.TypeOf(routeInfo),
		routeInfo.GetMethod(),
		routeInfo.GetPath(),
	)
	router.router.Handle(routeInfo.GetMethod(), routeInfo.GetPath(), router.wrap(handler))

	return nil
}

func (router *_router) wrap(handle route.Handler) httprouter.Handle {
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
