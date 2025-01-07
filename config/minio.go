package config

import "os"

type MinioConfig struct {
	Endpoint             string
	Host                 string
	Port                 string
	User                 string
	Password             string
	BucketProfile        string
	BucketStudentCode    string
	BucketSupervisorCode string
}

func LoadMinioConfig() MinioConfig {
	return MinioConfig{
		Endpoint:             os.Getenv("MINIO_ENDPOINT"),
		Host:                 os.Getenv("MINIO_HOST"),
		Port:                 os.Getenv("MINIO_PORT"),
		User:                 os.Getenv("MINIO_ROOT_USER"),
		Password:             os.Getenv("MINIO_ROOT_PASSWORD"),
		BucketProfile:        os.Getenv("MINIO_BUCKET_PROFILE"),
		BucketStudentCode:    os.Getenv("MINIO_BUCKET_STUDENT_CODE"),
		BucketSupervisorCode: os.Getenv("MINIO_BUCKET_SUPERVISOR_CODE"),
	}
}
