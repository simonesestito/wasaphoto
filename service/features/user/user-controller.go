package user

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

type Controller struct {
	// Dependencies
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  "GET",
			Path:    "/users/:userId",
			Handler: controller.getUserProfile,
		},
	}
}

func (controller Controller) getUserProfile(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	foundUser, err := controller.Service.GetUserAs(args.UserId, context.UserId)
	switch {
	case err == ErrWrongUUID:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case err == ErrUserBanned:
		http.Error(w, err.Error(), http.StatusForbidden)
	case err != nil:
		context.Logger.WithError(err).Error("error getting user as me")
		http.Error(w, "", http.StatusInternalServerError)
	case foundUser == nil:
		http.Error(w, "", http.StatusNotFound)
	default:
		api.SendJson(w, foundUser, 200, context.Logger)
	}
}
