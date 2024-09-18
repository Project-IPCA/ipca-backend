package builders

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type GroupChapterSelectedItemBuilder struct {
	GroupID    uuid.UUID
	ChapterID  uuid.UUID
	ItemID     int64
	ExerciseID uuid.UUID
}

func NewGroupChapterSelectedItemBuilder() *GroupChapterSelectedItemBuilder {
	return &GroupChapterSelectedItemBuilder{}
}

func (builder *GroupChapterSelectedItemBuilder) SetGroupID(groupID uuid.UUID) *GroupChapterSelectedItemBuilder {
	builder.GroupID = groupID
	return builder
}

func (builder *GroupChapterSelectedItemBuilder) SetChapterID(chapterID uuid.UUID) *GroupChapterSelectedItemBuilder {
	builder.ChapterID = chapterID
	return builder
}

func (builder *GroupChapterSelectedItemBuilder) SetItemID(itemID int64) *GroupChapterSelectedItemBuilder {
	builder.ItemID = itemID
	return builder
}

func (builder *GroupChapterSelectedItemBuilder) SetExerciseID(exerciseID uuid.UUID) *GroupChapterSelectedItemBuilder {
	builder.ExerciseID = exerciseID
	return builder
}

func (builder *GroupChapterSelectedItemBuilder) Build() models.GroupChapterSelectedItem {
	return models.GroupChapterSelectedItem{
		GroupID:    builder.GroupID,
		ChapterID:  builder.ChapterID,
		ItemID:     builder.ItemID,
		ExerciseID: builder.ExerciseID,
	}
}