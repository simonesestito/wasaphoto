package utils

import (
	"errors"
	"github.com/julienschmidt/httprouter"
)

// ParamsToMap converts httprouter.Params to a simple map[string]string
// to have a nice and simple interface,
// without the implications of a library
func ParamsToMap(params httprouter.Params) (map[string]string, error) {
	resultingMap := make(map[string]string, len(params))

	for _, entry := range params {
		key, value := entry.Key, entry.Value

		_, alreadyPresent := resultingMap[key]
		if alreadyPresent {
			return resultingMap, errors.New("Multiple values found for key " + key)
		}

		resultingMap[key] = value
	}

	return resultingMap, nil
}
