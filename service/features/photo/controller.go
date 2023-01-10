package photo

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"io"
	"net/http"
)

type Controller struct {
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  http.MethodPost,
			Path:    "/photos/",
			Handler: controller.uploadPhoto,
		},
		route.SecureRoute{
			Method:  http.MethodDelete,
			Path:    "/photos/:photoId",
			Handler: controller.deletePhoto,
		},
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId/photos/",
			Handler: controller.listUserPhotos,
		},
	}
}

func (controller Controller) uploadPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, context route.SecureRequestContext) {
	// Read photo file from body
	photoData, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		context.Logger.WithError(err).Errorln("error receiving photo")
		http.Error(w, "unexpected error receiving photo", http.StatusInternalServerError)
		return
	}

	if len(photoData) == 0 {
		http.Error(w, "missing photo body", http.StatusBadRequest)
		return
	}

	photo, err := controller.Service.CreatePost(context.UserId, photoData, context.Logger)

	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusCreated, context.Logger)
	} else {
		photo.AddImageHost(r, context.Logger)
		api.SendJson(w, photo, http.StatusCreated, context.Logger)
	}
}

func (controller Controller) deletePhoto(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParam{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	err := controller.Service.DeletePostAs(args.PhotoId, context.UserId)
	if errors.Is(err, api.ErrNotFound) {
		// Ignore, since it's intended to be idempotent
		err = nil
	}
	api.HandleErrorsResponse(err, w, http.StatusNoContent, context.Logger)
}

func (controller Controller) listUserPhotos(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseAllRequestVariables(r, params, &UserPhotosCursor{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	photos, cursor, err := controller.Service.GetUsersPhotosPage(args.UserId, context.UserId, args.PageCursorOrEmpty)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		// Add photo URL prefix
		for i := range photos {
			photos[i].AddImageHost(r, context.Logger)
		}

		api.SendJson(w, api.PageResult[Photo]{
			NextPageCursor: cursor,
			PageData:       photos,
		}, http.StatusOK, context.Logger)
	}
}
