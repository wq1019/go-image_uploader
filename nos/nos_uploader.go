package nos

import (
	"fmt"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/model"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosconst"
	. "github.com/wq1019/go-image_uploader"
	"io"
	"mime"
	"path/filepath"
)

type nosUploader struct {
	h          Hasher
	s          Store
	client     *nosclient.NosClient
	bucketName string
	h2sn       Hash2StorageName
}

func (nu *nosUploader) saveToNos(hashValue string, fh FileHeader, info ImageInfo) error {
	name, err := nu.h2sn.Convent(hashValue)
	if err != nil {
		return err
	}
	_, err = fh.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// 在 apline 镜像中 mime.TypeByExtension 只能用 jpg
	if info.Format == "jpeg" {
		info.Format = "jpg"
	}

	// Nos 只允许 最大 100MB 的文件
	_, err = nu.client.PutObjectByStream(&model.PutObjectRequest{
		Bucket: nu.bucketName,
		Object: name,
		Body:   fh.File,
		Metadata: &model.ObjectMetadata{
			ContentLength: fh.Size,
			Metadata: map[string]string{
				nosconst.CONTENT_TYPE: mime.TypeByExtension("." + info.Format),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("nos client put object stream error. err: %+v", err)
	}
	return err
}

func (nu *nosUploader) Upload(fh FileHeader) (*Image, error) {
	info, err := DecodeImageInfo(fh)
	if err != nil {
		return nil, err
	}
	hashValue, err := nu.h.Hash(fh.File)
	if err != nil {
		return nil, err
	}
	if exist, err := nu.s.ImageExist(hashValue); exist && err == nil {
		// 图片已经存在
		return nu.s.ImageLoad(hashValue)
	} else if err != nil {
		return nil, err
	}

	if err := nu.saveToNos(hashValue, fh, info); err != nil {
		return nil, err
	}

	return SaveToStore(nu.s, hashValue, fh.Filename, info)
}

func (nu *nosUploader) UploadFromURL(u string, filename string) (*Image, error) {
	if filename == "" {
		filename = filepath.Base(u)
	}
	file, size, err := DownloadImage(u)

	if err != nil {
		return nil, err
	}

	defer RemoveFile(file)

	fh := FileHeader{
		Filename: filename,
		Size:     size,
		File:     file,
	}

	return nu.Upload(fh)
}

func NewNosUploader(h Hasher, s Store, nosClient *nosclient.NosClient, bucketName string, h2sn Hash2StorageName) Uploader {
	if h2sn == nil {
		h2sn = Hash2StorageNameFunc(DefaultHash2StorageNameFunc)
	}
	return &nosUploader{h: h, s: s, client: nosClient, bucketName: bucketName, h2sn: h2sn}
}
