package route

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler = func(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
	context RequestContext,
)

type SecureHandler func(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
	context SecureRequestContext,
)

type Middleware = func(handler Handler) Handler

type AuthMiddleware interface {
	Intercept(handler SecureHandler) Handler
}
