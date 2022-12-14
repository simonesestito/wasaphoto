package api

import (
	"errors"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleErrorsResponse(err error, w http.ResponseWriter, defaultSuccessStatus int, logger logrus.FieldLogger) {
	switch {
	case errors.Is(err, ErrWrongUUID):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrWrongCursor):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrSelfOperation):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrDuplicated):
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, ErrAlreadyTaken):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, ErrUserBanned):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, ErrMedia):
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	case errors.Is(err, ErrOthersData):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, ErrThirdParty):
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case errors.Is(err, nil):
		w.WriteHeader(defaultSuccessStatus)
	default:
		logger.WithError(err).Error("unexpected error processing request")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HandlePutResult handles the error ErrDuplicated in an idempotent PUT operation.
func HandlePutResult(result any, err error, w http.ResponseWriter, logger logrus.FieldLogger) {
	switch {
	case err == nil:
		// Success, new item has been created!
		SendJson(w, result, http.StatusCreated, logger)
	case errors.Is(err, ErrDuplicated):
		// Ignore duplicated errors, since it's idempotent.
		SendJson(w, result, http.StatusOK, logger)
	default:
		HandleErrorsResponse(err, w, http.StatusOK, logger)
	}
}

// ErrSelfOperation is used to indicate a user is performing
// an operation both as subject and object,
// and that is not possible in this circumstance.
var ErrSelfOperation = errors.New("operation not allowed on yourself")

// ErrNotFound is used if the object of an operation cannot be found
var ErrNotFound = errors.New("subject resource not found")

// ErrDuplicated is used if an item was already added in a set
var ErrDuplicated = database.ErrDuplicated

// ErrAlreadyTaken is used if a user tries to get something it's someone else's and must be unique
var ErrAlreadyTaken = errors.New("the unique data you are trying to get it's already taken")

// ErrWrongUUID is used to indicate the given ID cannot be interpreted as a UUID
var ErrWrongUUID = errors.New("wrong UUID supplied")

// ErrWrongCursor is used to indicate the given pageCursor cannot be interpreted
var ErrWrongCursor = errors.New("wrong page cursor format")

// ErrUserBanned is used in case the current user has no permission
// to read the requested information because he is banned
// by the owner of that data.
var ErrUserBanned = errors.New("forbidden because of user ban")

// ErrMedia indicates a wrong media type
var ErrMedia = errors.New("wrong media content supplied")

// ErrOthersData indicates the current operation
// can only operate on data owned by the user performing it
var ErrOthersData = errors.New("you are only allowed to operate on data of yours")

// ErrThirdParty indicates the current operation cannot be fulfilled because,
// even if this server works as expected,
// it needs a third-party service which is currently not working.
var ErrThirdParty = errors.New("a third-party service required to fulfill the request is not available")
