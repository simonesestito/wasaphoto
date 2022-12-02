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

func MapGetFirstValue[T any](multiValueMap map[string][]T) (map[string]T, error) {
	singleValueMap := make(map[string]T)
	for key, values := range multiValueMap {
		if len(values) == 1 {
			singleValueMap[key] = values[0]
		} else if len(values) > 1 {
			return singleValueMap, errors.New("Multiple values for key " + key)
		}
	}
	return singleValueMap, nil
}

func JoinMaps[T any](maps ...map[string]T) (map[string]T, error) {
	joinMap := make(map[string]T)
	for _, singleMap := range maps {
		for key, value := range singleMap {
			if _, keyAlreadyPresent := joinMap[key]; keyAlreadyPresent {
				// Duplicate key!
				return joinMap, errors.New("Multiple key found: " + key)
			}
			joinMap[key] = value
		}
	}
	return joinMap, nil
}
