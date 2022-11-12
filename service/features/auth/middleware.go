package auth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
	"strings"
)

type Middleware struct {
	LoginService LoginService
}

func (middleware Middleware) Intercept(handler route.SecureHandler) route.Handler {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context route.RequestContext) {
		const bearerPrefix = "Bearer "
		authorization := request.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, bearerPrefix) {
			// Not authenticated
			http.Error(writer, "Unauthorized", 401)
			return
		}

		authToken := strings.TrimPrefix(authorization, bearerPrefix)
		userId, err := middleware.LoginService.IsAuthenticated(authToken)
		if err != nil {
			// Authentication is invalid
			context.Logger.WithError(err).Debug("Error checking authentication")
			http.Error(writer, "Unauthorized", 401)
			return
		}

		// Create secure context
		secureContext := route.SecureRequestContext{
			RequestContext: context,
			UserId:         userId,
		}

		// Authentication is valid!
		handler(writer, request, params, secureContext)
	}
}
