package photo

import (
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
			Path:    "/photos",
			Handler: controller.uploadPhoto,
		},
	}
}

func (controller Controller) uploadPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, context route.SecureRequestContext) {
	// Read photo file from body
	photoData, err := io.ReadAll(r.Body)
	if err != nil {
		context.Logger.WithError(err).Errorln("error receiving photo")
		http.Error(w, "unexpected error receiving photo", http.StatusInternalServerError)
		return
	}

	if len(photoData) == 0 {
		http.Error(w, "missing photo body", http.StatusBadRequest)
		return
	}

	photo, err := controller.Service.CreatePost(context.UserId, photoData)

	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusCreated, context.Logger)
	} else {
		api.SendJson(w, photo, http.StatusCreated, context.Logger)
	}
}
