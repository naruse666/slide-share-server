package storage

import (
	"context"

	"firebase.google.com/go/v4/storage"
)

type IThumbnailRepository interface {
	UploadThumbnail(thumbnailData []byte, fileName, folderName string) (string, error)
}

type ThumbnailRepository struct {
	client *storage.Client
}

func NewThumbnailRepository(client *storage.Client) IThumbnailRepository {
	return &ThumbnailRepository{client: client}
}

func (tr *ThumbnailRepository) UploadThumbnail(thumbnailData []byte, fileName, folderName string) (string, error) {
	ctx := context.Background()
	bucket, err := tr.client.DefaultBucket()
	if err != nil {
		return "", err
	}

	object := bucket.Object(folderName + "/" + fileName)
	wc := object.NewWriter(ctx)
	wc.ContentType = "image/png"
	if _, err = wc.Write([]byte(thumbnailData)); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	publicURL := "https://storage.googleapis.com/" + attrs.Bucket + "/" + attrs.Name

	return publicURL, nil
}
