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

	switch err {
	case ErrWrongUUID:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrSelfOperation:
		http.Error(w, err.Error(), http.StatusConflict)
	case ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrDuplicated:
		w.WriteHeader(http.StatusNoContent)
	case nil:
		w.WriteHeader(http.StatusCreated)
	default:
		context.Logger.WithError(err).Error("unexpected ban error")
		w.WriteHeader(http.StatusInternalServerError)
	}
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

	switch err {
	case ErrWrongUUID:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		context.Logger.WithError(err).Error("unexpected ban error")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
