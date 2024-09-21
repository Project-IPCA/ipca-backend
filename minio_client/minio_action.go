package minioclient

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/Project-IPCA/ipca-backend/pkg/utils"
)

type MinioAction struct {
	Minio *minio.Client
}

func NewMinioAction(minio *minio.Client) *MinioAction {
	return &MinioAction{Minio: minio}
}

func (minioAction *MinioAction) UploadToMinio(
	file *multipart.FileHeader,
	bucketName string,
) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	objectUlid := utils.NewULID()
	objectName := fmt.Sprintf("%s%s", objectUlid, fileExtension)

	contentType := file.Header.Get("Content-Type")
	_, err = minioAction.Minio.PutObject(
		context.Background(),
		bucketName,
		objectName,
		fileContent,
		file.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}

	return objectName, nil
}
