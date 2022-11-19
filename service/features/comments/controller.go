package comments

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
			Method:  http.MethodPost,
			Path:    "/photos/:photoId/comments/",
			Handler: controller.commentPhoto,
		},
		route.SecureRoute{
			Method:  http.MethodDelete,
			Path:    "/photos/:photoId/comments/:commentId",
			Handler: controller.uncommentPhoto,
		},
	}
}

func (controller Controller) commentPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, body, bodyErr := api.ParseVariablesAndBody(r, params, &photo.IdParam{}, &NewComment{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	createdComment, err := controller.Service.CommentPhoto(args.PhotoId, context.UserId, *body)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusCreated, context.Logger)
		return
	}

	api.SendJson(w, createdComment, http.StatusCreated, context.Logger)
}

func (controller Controller) uncommentPhoto(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	err := controller.Service.DeleteCommentOnPhotoIfAuthor(args.CommentId, args.PhotoId, context.UserId)
	api.HandleErrorsResponse(err, w, http.StatusNoContent, context.Logger)
}
