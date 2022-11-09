package auth

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/route"
	"github.com/julienschmidt/httprouter"
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
			Path:    "/sessions/",
			Handler: controller.handlePost,
		},
	}
}

func (controller LoginController) handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx route.RequestContext) {
	body, bodyErr := api.ParseAndValidateBody(r, &UserLoginCredentials{}, ctx.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	authToken, err := controller.AuthService.Authenticate(*body)
	if err != nil {
		ctx.Logger.WithError(err).Warnf("Unexpected error authenticating user %s", body.Username)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserLoginResult{UserId: authToken}
	api.SendJson(w, response, 200, ctx.Logger)
}
