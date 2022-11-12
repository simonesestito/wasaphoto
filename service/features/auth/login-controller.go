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
	responseStatus := http.StatusOK

	switch {
	case err == ErrUnknownUser:
		// Register new user!
		responseStatus = http.StatusCreated
		authToken, err = controller.AuthService.SignUp(*body)
		if err != nil {
			ctx.Logger.WithError(err).Errorf("Unexpected error signing up user with username '%s'", body.Username)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case err != nil:
		ctx.Logger.WithError(err).Warnf("Unexpected error authenticating user %s", body.Username)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserLoginResult{UserId: authToken}
	api.SendJson(w, response, responseStatus, ctx.Logger)
}
