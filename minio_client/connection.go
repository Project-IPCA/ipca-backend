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

	bucketNames := []string{cfg.Minio.BucketProfile, cfg.Minio.BucketStudentCode, cfg.Minio.BucketSupervisorCode}
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

			readOnlyPolicy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"AWS": ["*"]
						},
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::%s/*"]
					}
				]
			}`, bucketName)
		
			// Apply the policy to the bucket
			err = minioClient.SetBucketPolicy(context.Background(), bucketName, readOnlyPolicy)
			if err != nil {
				panic("Failed to set bucket policy:"+ err.Error())
			}
		}
	}

	return minioClient
}
