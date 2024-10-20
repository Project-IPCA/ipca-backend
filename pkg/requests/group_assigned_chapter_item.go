package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UpdateGroupAssignedChapterItemRequest struct{
	ItemId	int `json:"item_id" validate:"required"       example:"1"`
	ChapterId uuid.UUID `json:"chapter_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
	GroupId	uuid.UUID `json:"group_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
	SelectedItem  []uuid.UUID `json:"select_items"      validate:"required"`
}

func (asr UpdateGroupAssignedChapterItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&asr)
}

type UpdateAllGroupAssignedChapterItemRequest struct{
	ChapterId uuid.UUID `json:"chapter_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
	GroupId	uuid.UUID `json:"group_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
	UpdatePool []AllUpdateData `json:"update_pool" validate:"required"`
}

type AllUpdateData struct {
	ItemId	int `json:"item_id" validate:"required"       example:"1"`
	SelectedItem  []uuid.UUID `json:"select_items"      validate:"required"`
}

func (asr UpdateAllGroupAssignedChapterItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&asr)
}