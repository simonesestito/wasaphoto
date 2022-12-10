package auth

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

type LoginController struct {
	// Dependencies
	AuthService LoginService
}

func (controller LoginController) ListRoutes() []route.Route {
	return []route.Route{
		route.AnonymousRoute{
			Method:  http.MethodPost,
			Path:    "/session",
			Handler: controller.handlePost,
		},
	}
}

func (controller LoginController) handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx route.RequestContext) {
	body, bodyErr := api.ParseAndValidateBody(r, &userLoginCredentials{}, ctx.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	authToken, isNew, err := controller.AuthService.AuthenticateOrSignup(*body)

	var responseStatus int
	var logAction string
	if isNew {
		responseStatus = http.StatusCreated
		logAction = "signing up"
	} else {
		responseStatus = http.StatusOK
		logAction = "authenticating"
	}

	if err != nil {
		ctx.Logger.WithError(err).Warnf("Unexpected error %s user with username '%s'", logAction, body.Username)
		api.HandleErrorsResponse(err, w, responseStatus, ctx.Logger)
		return
	}

	response := userLoginResult{UserId: authToken}
	api.SendJson(w, response, responseStatus, ctx.Logger)
}
