package tinypng

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

// In a real project, always use config file for that (not git committed)
const apiToken = "lHLbzbX5M7v4NqtPVPzjqqqY4zJvHgy8"

type API struct {
}

func (imgApi API) CompressPhoto(imageData []byte, logger logrus.FieldLogger) ([]byte, error) {
	logger.Debugln("Uploading photo...")
	imageId, err := imgApi.uploadPhoto(imageData)
	if err != nil {
		return nil, err
	}

	logger.Debugln("Converting to WebP...")
	webp, err := imgApi.convertToWebp(imageId)
	if err != nil {
		logger.WithError(err).Errorln("Conversion to WebP failed")
		return nil, err
	}

	logger.Debugln("Conversion to WebP successful")
	return webp, nil
}

func (API) getAuthString() string {
	authString := "api:" + apiToken
	authDigest := base64.StdEncoding.EncodeToString([]byte(authString))
	return "Basic " + authDigest
}

func (imgApi API) uploadPhoto(imageData []byte) (string, error) {
	apiUrl, err := url.Parse("https://api.tinify.com/shrink")
	if err != nil {
		return "", err
	}

	request := &http.Request{
		Method: http.MethodPost,
		URL:    apiUrl,
		Header: map[string][]string{
			// "Content-Type":  {"image/*"},
			"Authorization": {imgApi.getAuthString()},
		},
		Body:          io.NopCloser(bytes.NewReader(imageData)),
		ContentLength: int64(len(imageData)),
		Host:          "api.tinify.com",
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusCreated:
		return response.Header.Get("Location"), nil
	case http.StatusUnsupportedMediaType:
		return "", api.ErrMedia
	default:
		return "", errors.New("TinyPNG API status code response: " + response.Status)
	}
}

func (imgApi API) convertToWebp(imageId string) ([]byte, error) {
	apiUrl, err := url.Parse(imageId)
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(map[string]any{
		"convert": map[string]string{"type": "image/webp"},
		"resize": map[string]any{
			"method": "scale",
			"width":  500,
		},
	})
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: http.MethodPost,
		URL:    apiUrl,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {imgApi.getAuthString()},
		},
		Body:          io.NopCloser(bytes.NewReader(jsonBody)),
		ContentLength: int64(len(jsonBody)),
		Host:          "api.tinify.com",
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("TinyPNG API status code response: " + response.Status)
	}

	// Read WebP file
	return io.ReadAll(response.Body)
}
