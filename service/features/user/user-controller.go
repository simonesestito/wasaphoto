package user

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/api/route"
	"io"
	"net/http"
)

type Controller struct {
	// Dependencies
	Service Service
}

func (controller Controller) ListRoutes() []route.Route {
	return []route.Route{
		route.SecureRoute{
			Method:  http.MethodGet,
			Path:    "/users/:userId",
			Handler: controller.getUserProfile,
		},
		route.SecureRoute{
			Method:  http.MethodPut,
			Path:    "/users/:userId",
			Handler: controller.setMyDetails,
		},
		route.SecureRoute{
			Method:  http.MethodPut,
			Path:    "/users/:userId/username",
			Handler: controller.setMyUserName,
		},
	}
}

func (controller Controller) getUserProfile(w http.ResponseWriter, _ *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	foundUser, err := controller.Service.GetUserAs(args.UserId, context.UserId)
	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else if foundUser == nil {
		http.Error(w, "not found", http.StatusNotFound)
	} else {
		api.SendJson(w, foundUser, 200, context.Logger)
	}
}

func (controller Controller) setMyDetails(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	body, bodyErr := api.ParseAndValidateBody(r, &NewUser{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	updatedUser, err := controller.Service.UpdateUserDetails(args.UserId, *body)

	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		api.SendJson(w, updatedUser, 200, context.Logger)
	}
}

func (controller Controller) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context route.SecureRequestContext) {
	// Validate path arguments
	args, bodyErr := api.ParseRequestVariables(params, &IdParams{}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	if args.UserId != context.UserId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Read username from request
	usernameBytes, err := io.ReadAll(r.Body)
	if err != nil {
		context.Logger.WithError(err).Debugln("Error reading username from body")
		http.Error(w, "unable to read username", http.StatusBadRequest)
		return
	}
	username := string(usernameBytes)

	// Validate read username
	bodyErr = api.ValidateParsedStruct(&UsernameGetParams{username}, context.Logger)
	if bodyErr != nil {
		http.Error(w, bodyErr.Message, bodyErr.StatusCode)
		return
	}

	updatedUser, err := controller.Service.UpdateUsername(args.UserId, username)

	if err != nil {
		api.HandleErrorsResponse(err, w, http.StatusOK, context.Logger)
	} else {
		api.SendJson(w, updatedUser, 200, context.Logger)
	}
}
