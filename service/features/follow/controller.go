package follow

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/user"
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
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId/followers/",
			Handler: controller.listFollowers,
		},
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId/followings/",
			Handler: controller.listFollowings,
		},
	}
}

func (controller Controller) followUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &followerParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.FollowUser(context.UserId, args.FollowedId)
	result := userFollow{
		FollowingId: args.FollowedId,
		FollowerId:  context.UserId,
	}
	api.HandlePutResult(result, err, w, context.Logger)
}

func (controller Controller) unfollowUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &followerParams{}, context.Logger)
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

func (controller Controller) listFollowers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseAllRequestVariables(r, params, &user.IdUserCursor{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	followers, cursor, err := controller.Service.ListFollowersAs(args.UserId, context.UserId, args.PageCursorOrEmpty)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		api.SendJson(w, api.PageResult[user.User]{
			NextPageCursor: cursor,
			PageData:       followers,
		}, http.StatusOK, context.Logger)
	}
}
func (controller Controller) listFollowings(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseAllRequestVariables(r, params, &user.IdUserCursor{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	followers, cursor, err := controller.Service.ListFollowingsAs(args.UserId, context.UserId, args.PageCursorOrEmpty)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		api.SendJson(w, api.PageResult[user.User]{
			NextPageCursor: cursor,
			PageData:       followers,
		}, http.StatusOK, context.Logger)
	}
}
