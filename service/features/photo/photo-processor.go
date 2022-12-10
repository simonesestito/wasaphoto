package photo

import "github.com/simonesestito/wasaphoto/service/utils/tinypng"

type imageProcessor struct {
	TinyPng tinypng.API
}

func (processor imageProcessor) compressPhotoToWebp(imageData []byte) ([]byte, error) {
	return processor.TinyPng.CompressPhoto(imageData)
}
