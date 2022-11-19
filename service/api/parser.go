package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	"github.com/simonesestito/wasaphoto/service/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type MalformedRequest struct {
	StatusCode int
	Message    string
}

func (err *MalformedRequest) Error() string { return err.Message }

func ParseRequestVariables[T any](params httprouter.Params, paramsStruct *T, logger logrus.FieldLogger) (*T, *MalformedRequest) {
	paramsMap, err := utils.ParamsToMap(params)
	if err != nil {
		return nil, &MalformedRequest{http.StatusBadRequest, err.Error()}
	}

	// Convert allParamsMap to a struct
	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		TagName:     "json",
		Squash:      true,
		Result:      paramsStruct,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		msg := "Error creating a decoder"
		logger.Errorf("%s: %s", msg, err.Error())
		return nil, &MalformedRequest{http.StatusInternalServerError, msg}
	}

	if err := decoder.Decode(paramsMap); err != nil {
		return nil, &MalformedRequest{http.StatusBadRequest, err.Error()}
	}

	// Validate parsed struct
	if err := ValidateParsedStruct(paramsStruct, logger); err != nil {
		return paramsStruct, err
	}

	return paramsStruct, nil
}

func ParseAndValidateBody[T any](request *http.Request, bodyStruct *T, logger logrus.FieldLogger) (*T, *MalformedRequest) {
	// Close request body at the end of parsing
	defer request.Body.Close()

	// Check content type
	if request.Header.Get("Content-Type") != "application/json" {
		return nil, &MalformedRequest{http.StatusUnsupportedMediaType, "Content-Type is not JSON"}
	}

	// Decode JSON body
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(bodyStruct); err != nil {
		return nil, &MalformedRequest{http.StatusBadRequest, explainJsonError(err, logger)}
	}

	// Check if no more JSON is available
	if decoder.Decode(&struct{}{}) != io.EOF {
		return nil, &MalformedRequest{
			http.StatusBadRequest, "Request body can only contain one JSON object",
		}
	}

	// Validate parsed struct
	if err := ValidateParsedStruct(bodyStruct, logger); err != nil {
		return bodyStruct, err
	}

	return bodyStruct, nil
}

func ValidateParsedStruct[T any](parsedStruct *T, logger logrus.FieldLogger) *MalformedRequest {
	if err := validator.New().Struct(parsedStruct); err != nil {
		validationError, ok := err.(validator.ValidationErrors)
		if !ok {
			msg := "Unexpected error validating body input"
			logger.Errorf("%s: %s", msg, err.Error())
			return &MalformedRequest{http.StatusInternalServerError, msg}
		}

		msg := "Error validating input body"
		if len(validationError) > 0 {
			err := validationError[0]
			msg = fmt.Sprintf("%s: %s %s %s", msg, err.Namespace(), err.Tag(), err.Param())
		}

		return &MalformedRequest{http.StatusBadRequest, msg}
	}

	return nil
}

func explainJsonError(err error, logger logrus.FieldLogger) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	// From: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
	switch {
	case errors.As(err, &syntaxError):
		return fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Sprintf("Request body contains badly-formed JSON")

	case errors.As(err, &unmarshalTypeError):
		return fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Sprintf("Request body contains unknown field %s", fieldName)

	case errors.Is(err, io.EOF):
		return "Request body must not be empty"

	default:
		msg := "Unexpected error parsing JSON"
		logger.WithError(err).Warnf(msg)
		return msg
	}
}

func ParseVariablesAndBody[V any, B any](r *http.Request, params httprouter.Params, varsStruct *V, bodyStruct *B, logger logrus.FieldLogger) (*V, *B, *MalformedRequest) {
	args, err := ParseRequestVariables(params, varsStruct, logger)
	if err != nil {
		return nil, nil, err
	}

	body, err := ParseAndValidateBody(r, bodyStruct, logger)
	if err != nil {
		return nil, nil, err
	}

	return args, body, nil
}

func SendJson(writer http.ResponseWriter, body any, successStatus int, logger logrus.FieldLogger) {
	response, err := json.Marshal(body)
	if err != nil {
		logger.WithError(err).Error("Error marshalling JSON response")
		http.Error(writer, "Internal Server Error", 500)
	} else {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(successStatus)
		_, err := writer.Write(response)
		if err != nil {
			logger.WithError(err).Warn("Error sending response data to client")
		}
	}
}
