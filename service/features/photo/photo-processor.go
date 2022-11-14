package photo

type ImageProcessor struct {
}

func (ImageProcessor) CompressPhotoToWebp(imageData []byte) ([]byte, error) {
	return imageData, nil
	// TODO: Handle non image file error (api.ErrMedia)
}
