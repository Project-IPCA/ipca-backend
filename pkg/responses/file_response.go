package responses

type FileResponse struct {
	ObjectName string `json:"object_name"`
	ObjectUrl  string `json:"object_url"`
}

func NewFileResponse(objectName string, objectUrl string) *FileResponse {
	return &FileResponse{
		ObjectName: objectName,
		ObjectUrl:  objectUrl,
	}
}
