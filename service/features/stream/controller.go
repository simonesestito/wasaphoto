package stream

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"github.com/simonesestito/wasaphoto/service/features/photo"
	"net/http"
)

type Controller struct {
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId/stream",
			Handler: controller.getMyStream,
		},
	}
}

func (controller Controller) getMyStream(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseAllRequestVariables(r, params, &photo.UserPhotosCursor{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	nextPhotos, nextCursor, err := controller.Service.GetStreamPage(context.UserId, args.PageCursorOrEmpty)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		// Add photo URL prefix
		for i := range nextPhotos {
			nextPhotos[i].AddImageHost(r, context.Logger)
		}

		page := api.PageResult[photo.Photo]{
			NextPageCursor: nextCursor,
			PageData:       nextPhotos,
		}
		api.SendJson(w, page, http.StatusOK, context.Logger)
	}
}
