package user

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

type BanController struct {
	Service BanService
}

func (controller BanController) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  http.MethodPut,
			Path:    "/users/:userId/bannedPeople/:bannedId",
			Handler: controller.banUser,
		},
		route.SecureRoute{
			Method:  http.MethodDelete,
			Path:    "/users/:userId/bannedPeople/:bannedId",
			Handler: controller.unbanUser,
		},
	}
}

func (controller BanController) banUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &banParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.BanUser(args.BannedId, args.UserId)
	result := banResult{
		BannedId: args.BannedId,
		BannerId: context.UserId,
	}
	api.HandlePutResult(result, err, w, context.Logger)
}

func (controller BanController) unbanUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &banParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.UnbanUser(args.BannedId, args.UserId)
	api.HandleErrorsResponse(err, w, http.StatusNoContent, context.Logger)
}
