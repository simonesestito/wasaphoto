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
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId/bannedPeople/",
			Handler: controller.listBannedUsers,
		},
	}
}

func (controller BanController) banUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &BanParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.BanUser(args.BannedId, args.UserId)
	result := BanResult{
		BannedId: args.BannedId,
		BannerId: context.UserId,
	}
	api.HandlePutResult(result, err, w, context.Logger)
}

func (controller BanController) unbanUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &BanParams{}, context.Logger)
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

func (controller BanController) listBannedUsers(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	users, err := controller.Service.GetBannedUsers(context.UserId)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		api.SendJson(w, users, http.StatusOK, context.Logger)
	}
}
