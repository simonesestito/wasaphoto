package api

import (
	"errors"
	"github.com/simonesestito/wasaphoto/service/database"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleErrorsResponse(err error, w http.ResponseWriter, defaultSuccessStatus int, logger logrus.FieldLogger) {
	switch err {
	case ErrWrongUUID:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrSelfOperation:
		http.Error(w, err.Error(), http.StatusConflict)
	case ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrDuplicated:
		w.WriteHeader(http.StatusNoContent)
	case ErrAlreadyTaken:
		http.Error(w, err.Error(), http.StatusConflict)
	case ErrUserBanned:
		http.Error(w, err.Error(), http.StatusForbidden)
	case ErrMedia:
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	case nil:
		w.WriteHeader(defaultSuccessStatus)
	default:
		logger.WithError(err).Error("unexpected error processing request")
		w.WriteHeader(http.StatusInternalServerError)
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

// ErrUserBanned is used in case the current user has no permission
// to read the requested information because he is banned
// by the owner of that data.
var ErrUserBanned = errors.New("forbidden because of user ban")

// ErrMedia indicates a wrong media type
var ErrMedia = errors.New("wrong media content supplied")
