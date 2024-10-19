package requests

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SetChapterPemissionRequest struct {
	Permission AccessAndSubmitData `json:"permission" validate:"required"`
	ChapterId uuid.UUID `json:"chapter_id" validate:"required"`
	GroupId uuid.UUID `json:"group_id" validate:"required"`
	Sync *bool `json:"sync" validate:"required"`
}

func (ba SetChapterPemissionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&ba)
}

type ChapterPermission struct {
	Access AccessAndSubmitData 
	Submit AccessAndSubmitData
	FullMark int
	Name string
	ChapterId uuid.UUID
	NoItems int
	Status  string
}

type AccessAndSubmitData struct {
	TimeStart *time.Time `json:"time_start"`
	TimeEnd		*time.Time `json:"time_end"`
	Type       string `json:"type" validate:"required"`
	Prefix string `json:"prefix" validate:"required"`
}