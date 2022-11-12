package follow

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"net/http"
)

type Controller struct {
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  http.MethodPut,
			Path:    "/users/:userId/followings/:followedId",
			Handler: controller.followUser,
		},
		route.SecureRoute{
			Method:  http.MethodDelete,
			Path:    "/users/:userId/followings/:followedId",
			Handler: controller.unfollowUser,
		},
	}
}

func (controller Controller) followUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &FollowerParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.FollowUser(context.UserId, args.FollowedId)
	api.HandleErrorsResponse(err, w, http.StatusCreated, context.Logger)
}

func (controller Controller) unfollowUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &FollowerParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.UnfollowUser(context.UserId, args.FollowedId)
	api.HandleErrorsResponse(err, w, http.StatusNoContent, context.Logger)
}
