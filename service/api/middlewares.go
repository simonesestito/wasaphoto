package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

// LimitBodySize prevents spending too much time
// responding to too long (potentially bad) requests
func LimitBodySize(maxBytes int64) route.Middleware {
	return func(handler route.Handler) route.Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.RequestContext) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			handler(w, r, params, context)
		}
	}
}
