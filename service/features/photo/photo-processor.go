package photo

import (
	"github.com/simonesestito/wasaphoto/service/utils/tinypng"
	"github.com/sirupsen/logrus"
)

type imageProcessor struct {
	TinyPng tinypng.API
}

func (processor imageProcessor) compressPhotoToWebp(imageData []byte, logger logrus.FieldLogger) ([]byte, error) {
	return processor.TinyPng.CompressPhoto(imageData, logger)
}
