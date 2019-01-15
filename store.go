package image_uploader

type Store interface {
	ImageExist(hash string) (bool, error)
	ImageLoad(hash string) (*Image, error)
	ImageCreate(image *Image) error
}

func SaveToStore(s Store, hashValue, title string, info ImageInfo) (imageModel *Image, err error) {
	imageModel = &Image{
		Hash:   hashValue,
		Size:   info.Size,
		Format: info.Format,
		Title:  title,
		Width:  info.Width,
		Height: info.Height,
	}
	err = s.ImageCreate(imageModel)
	return
}
