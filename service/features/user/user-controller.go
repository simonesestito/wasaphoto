package user

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/route"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Controller struct {
	// Dependencies
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  "GET",
			Path:    "/users/{userId}",
			Handler: controller.getUserById,
		},
	}
}

func (controller Controller) getUserById(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, err := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	foundUser := map[string]string{
		"id":         args.UserId,
		"searchedBy": context.UserId,
	}
	api.SendJson(w, foundUser, 200, context.Logger)
}
