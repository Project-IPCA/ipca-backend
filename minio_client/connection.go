package minioclient

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Project-IPCA/ipca-backend/config"
)

func Init(cfg *config.Config) *minio.Client {
	endpoint := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.User, cfg.Minio.Password, ""),
		Secure: false,
	})
	if err != nil {
		panic("failed to connect to minio: " + err.Error())
	}

	existBuckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		panic("failed to list bucket minio: " + err.Error())
	}

	bucketNames := []string{cfg.Minio.BucketProfile, cfg.Minio.BucketStudentCode}
	for _, bucketName := range bucketNames {
		isExist := false
		for _, bucket := range existBuckets {
			if bucket.Name == bucketName {
				isExist = true
				break
			}
		}

		if !isExist {
			err = minioClient.MakeBucket(
				context.Background(),
				bucketName,
				minio.MakeBucketOptions{},
			)
			if err != nil {
				panic("failed to create bucket minio: " + err.Error())
			}
		}
	}

	return minioClient
}
