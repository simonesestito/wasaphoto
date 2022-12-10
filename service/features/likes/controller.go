package likes

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
			Path:    "/photos/:photoId/likes/:userId",
			Handler: controller.likePhoto,
		},
		route.SecureRoute{
			Method:  http.MethodDelete,
			Path:    "/photos/:photoId/likes/:userId",
			Handler: controller.unlikePhoto,
		},
	}
}

func (controller Controller) likePhoto(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &likeParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.LikePhoto(args.PhotoId, context.UserId)
	result := photoLike{
		PhotoId: args.PhotoId,
		UserId:  args.UserId,
	}
	api.HandlePutResult(result, err, w, context.Logger)
}

func (controller Controller) unlikePhoto(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &likeParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := controller.Service.UnlikePhoto(args.PhotoId, context.UserId)
	api.HandleErrorsResponse(err, w, http.StatusNoContent, context.Logger)
}
