package image_uploader

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	// todo 下面这两个类型可以考虑不要
	//_ "golang.org/x/image/bmp"
	//_ "golang.org/x/image/tiff"
)

type ImageInfo struct {
	Width  uint
	Height uint
	Format string
	Size   int64
}

func DecodeImageInfo(fh FileHeader) (info ImageInfo, err error) {
	_, err = fh.File.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	config, format, err := image.DecodeConfig(fh.File)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		Width:  uint(config.Width),
		Height: uint(config.Height),
		Format: format,
		Size:   fh.Size,
	}, nil

}

func IsUnknownFormat(err error) bool {
	return image.ErrFormat == err
}
