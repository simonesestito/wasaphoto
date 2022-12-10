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

// RegisterAll registers all the received controllers, which in turn have a list of single API handlers.
// The controllers are declared and instantiated in ioc/controller.go.
//
// Each feature-specific controller has a `ListRoutes()` method where it declares all its handlers.
// Then, this function will register them all one by one.
//
// By doing so, the controller is not coupled with the registration mechanism and other middlewares.
// In fact, it's perfectly possible to use another routing method and still no need to touch controllers.
//
// For documentation purposes, here they are all:
// - User related endpoints are registered in features/user/user-controller.go (user.Controller#ListRoutes())
// -- 'route.SecureRoute' [GET] /users/:userId
// -- 'route.SecureRoute' [PUT] /users/:userId
// -- 'route.SecureRoute' [PUT] /users/:userId/username
// -- 'route.SecureRoute' [GET] /users/
//
// - Ban related endpoints are registered in features/user/ban-controller.go (user.BanController#ListRoutes())
// -- 'route.SecureRoute' [PUT] /users/:userId/bannedPeople/:bannedId
// -- 'route.SecureRoute' [DELETE] /users/:userId/bannedPeople/:bannedId
//
// - Stream related endpoints are registered in features/stream/controller.go (stream.Controller#ListRoutes())
// -- 'route.SecureRoute' [GET] /users/:userId/stream
//
// - Photo related endpoints are registered in features/photo/controller.go (photo.Controller#ListRoutes())
// -- 'route.SecureRoute' [POST] /photos
// -- 'route.SecureRoute' [DELETE] /photos/:photoId
// -- 'route.SecureRoute' [GET] /users/:userId/photos/
//
// - Likes related endpoints are registered in features/likes/controller.go (likes.Controller#ListRoutes())
// -- 'route.SecureRoute' [PUT] /photos/:photoId/likes/:userId
// -- 'route.SecureRoute' [DELETE] /photos/:photoId/likes/:userId
//
// - Follow related endpoints are registered in features/follow/controller.go (follow.Controller#ListRoutes())
// -- 'route.SecureRoute' [PUT] /users/:userId/followings/:followedId
// -- 'route.SecureRoute' [DELETE] /users/:userId/followings/:followedId
// -- 'route.SecureRoute' [GET] /users/:userId/followers/
// -- 'route.SecureRoute' [GET] /users/:userId/followings/
//
// - Comments related endpoints are registered in features/comments/controller.go (comments.Controller#ListRoutes())
// -- 'route.SecureRoute' [POST] /photos/:photoId/comments/
// -- 'route.SecureRoute' [DELETE] /photos/:photoId/comments/:commentId
// -- 'route.SecureRoute' [GET] /photos/:photoId/comments/
//
// - Login/auth related endpoints are registered in features/auth/login-controller.go (auth.LoginController#ListRoutes())
// -- 'route.AnonymousRoute' [POST] /session
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

// Handler returns the router http.Handler responsible for doing the routing of each request.
// Registration is performed in the RegisterAll function above.
func (router *_router) Handler() http.Handler {
	return router.router
}

// RegisterStatic handles requests for static files
func (router *_router) RegisterStatic(localPath string, httpPath string) {
	router.logger.Debugln("Registering static path", localPath, "to HTTP path", httpPath)
	router.router.ServeFiles(httpPath+"/*filepath", http.Dir(localPath))
}

// Register a new route.
// Normalize all different kinds of routes. For instance, a route.SecureRoute will
// be wrapped with a middleware which extracts auth info and performs some checks.
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
		return fmt.Errorf("unknown route type: %s", reflect.TypeOf(routeInfo))
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
