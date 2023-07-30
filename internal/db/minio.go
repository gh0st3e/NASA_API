package db

import (
	"context"
	"fmt"

	"github.com/gh0st3e/NASA_API/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func InitMinio(cfg config.MinioConfig) (*minio.Client, error) {
	ctx := context.Background()

	endpoint := cfg.Address
	accessKeyID := cfg.Username
	secretAccessKey := cfg.Password

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating minio client: %w", err)
	}

	err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.BucketName)
		if errBucketExists == nil && exists {
			logrus.Infof("bucket already exist %s\n", cfg.BucketName)
		} else {
			return nil, fmt.Errorf("error while creating bucket: %w", err)
		}
	}

	return minioClient, nil
}
