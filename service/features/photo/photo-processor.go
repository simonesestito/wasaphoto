package photo

import (
	"github.com/simonesestito/wasaphoto/service/api"
	"github.com/simonesestito/wasaphoto/service/utils/tinypng"
	"github.com/sirupsen/logrus"
)

type imageProcessor struct {
	TinyPng tinypng.API
}

func (processor imageProcessor) compressPhotoToWebp(imageData []byte, logger logrus.FieldLogger) ([]byte, error) {
	compressedImage, err := processor.TinyPng.CompressPhoto(imageData, logger)
	if err != nil {
		logger.WithError(err).Errorln("unable to reach TinyPNG service")
		return nil, api.ErrThirdParty
	}

	return compressedImage, nil
}
