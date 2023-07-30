package service

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

const (
	NoSuchKeyError = "NoSuchKey"
)

type ImageActions interface {
	PutImage(ctx context.Context, objectName, filePath, contentType string) (minio.UploadInfo, error)
	GetImage(ctx context.Context, objectName, filePath string) error
}

type ImageService struct {
	log        *logrus.Logger
	imageStore ImageActions
}

func NewImageService(log *logrus.Logger, imageStore ImageActions) *ImageService {
	return &ImageService{
		log:        log,
		imageStore: imageStore,
	}
}

func (i *ImageService) PutImage(ctx context.Context, objectName, filePath, contentType string) error {
	i.log.Info("[PutImage] started")

	info, err := i.imageStore.PutImage(ctx, objectName, filePath, contentType)
	if err != nil {
		i.log.Errorf("Error while putting image into storage: %s", err.Error())
		return fmt.Errorf("error while putting image into storage: %w", err)
	}

	i.log.Info(info)
	i.log.Info("[PutImage] ended")

	return nil
}

func (i *ImageService) GetImage(ctx context.Context, fileName string) (string, error) {
	i.log.Info("[GetImage] started")

	filePath := imagesPath + fileName + fileExtension

	err := i.imageStore.GetImage(ctx, fileName+fileExtension, filePath)
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok {
			if minioErr.Code == NoSuchKeyError {
				return "", fmt.Errorf("no such image")
			}
		}
		i.log.Errorf("Error while getting image from storage: %s", err.Error())
		return "", fmt.Errorf("error while getting image from storage: %w", err)
	}

	i.log.Info("[GetImage] ended")

	return filePath, nil
}
