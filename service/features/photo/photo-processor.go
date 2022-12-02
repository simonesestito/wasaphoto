package photo

import "github.com/simonesestito/wasaphoto/service/utils/tinypng"

type ImageProcessor struct {
	TinyPng tinypng.API
}

func (processor ImageProcessor) CompressPhotoToWebp(imageData []byte) ([]byte, error) {
	return processor.TinyPng.CompressPhoto(imageData)
}
