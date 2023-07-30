package store

import (
	"context"
	"fmt"

	minio "github.com/minio/minio-go/v7"
)

type ImageStore struct {
	minioClient *minio.Client
	bucketName  string
}

func NewImageStore(minioClient *minio.Client, bucketName string) *ImageStore {
	return &ImageStore{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (i *ImageStore) PutImage(ctx context.Context, objectName, filePath, contentType string) (minio.UploadInfo, error) {
	return i.minioClient.FPutObject(
		ctx,
		i.bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{ContentType: contentType})
}

func (i *ImageStore) GetImage(ctx context.Context, objectName, filePath string) error {
	fmt.Println(objectName)
	fmt.Println(filePath)
	return i.minioClient.FGetObject(
		ctx,
		i.bucketName,
		objectName,
		filePath,
		minio.GetObjectOptions{})
}
