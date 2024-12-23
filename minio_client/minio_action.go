package minioclient

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
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
	file interface{},
	bucketName string,
	isNotGenName bool,
) (string, error) {
	var fileContent io.Reader
	var fileName string
	var fileSize int64
	var contentType string

	switch f := file.(type) {
	case *multipart.FileHeader:
		var err error
		fileContent, err = f.Open()
		if err != nil {
			return "", fmt.Errorf("error opening multipart file: %w", err)
		}
		defer fileContent.(io.ReadCloser).Close()
		fileName = f.Filename
		fileSize = f.Size
		contentType = f.Header.Get("Content-Type")
	case *os.File:
		fileInfo, err := f.Stat()
		if err != nil {
			return "", fmt.Errorf("error getting file info: %w", err)
		}
		fileContent = f
		fmt.Println(fileInfo.Name())
		fileName = fileInfo.Name()
		fileSize = fileInfo.Size()
		contentType = "application/octet-stream"
	default:
		return "", fmt.Errorf("unsupported file type")
	}

	var objectName string
	if isNotGenName {
		objectName = fileName
	} else {
		fileExtension := strings.ToLower(filepath.Ext(fileName))
		objectUlid := utils.NewULID()
		objectName = fmt.Sprintf("%s%s", objectUlid, fileExtension)
	}

	_, err := minioAction.Minio.PutObject(
		context.Background(),
		bucketName,
		objectName,
		fileContent,
		fileSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error uploading to MinIO: %w", err)
	}

	return objectName, nil
}

func (minioAction *MinioAction) GetFromMinio(
	bucketName string,
	objectName string,
) (*minio.Object, error) {
	object, err := minioAction.Minio.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error while getting object from minio : %v", err)
	}
	return object, nil
}

func (minioAction *MinioAction) DeleteFileInMinio(
	bucketName string,
	objectName string,
) error {
	err := minioAction.Minio.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error while getting object from minio : %v", err)
	}
	return nil
}
