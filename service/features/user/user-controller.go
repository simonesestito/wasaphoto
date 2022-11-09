package user

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

type Controller struct {
	// Dependencies
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  "GET",
			Path:    "/users/:userId",
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
