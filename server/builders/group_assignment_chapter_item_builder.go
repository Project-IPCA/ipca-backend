package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupAssignmentChapterItemBuilder struct {
	GroupID        uuid.UUID
	ChapterID      uuid.UUID
	ItemID         int64
	ExerciseIDList *string
	FullMark       int
	TimeStart      *string
	TimeEnd        *string
	Status         *string
}

func NewGroupAssignmentChapterItemBuilder() *GroupAssignmentChapterItemBuilder {
	return &GroupAssignmentChapterItemBuilder{}
}

func (builder *GroupAssignmentChapterItemBuilder) SetGroupID(
	groupID uuid.UUID,
) *GroupAssignmentChapterItemBuilder {
	builder.GroupID = groupID
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetChapterID(
	chapterID uuid.UUID,
) *GroupAssignmentChapterItemBuilder {
	builder.ChapterID = chapterID
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetItemID(
	itemID int64,
) *GroupAssignmentChapterItemBuilder {
	builder.ItemID = itemID
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetExerciseIDList(
	exerciseIDList string,
) *GroupAssignmentChapterItemBuilder {
	builder.ExerciseIDList = &exerciseIDList
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetFullMark(
	fullMark int,
) *GroupAssignmentChapterItemBuilder {
	builder.FullMark = fullMark
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetTimeStart(
	timeStart *string,
) *GroupAssignmentChapterItemBuilder {
	builder.TimeStart = timeStart
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetTimeEnd(
	timeEnd *string,
) *GroupAssignmentChapterItemBuilder {
	builder.TimeEnd = timeEnd
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) SetStatus(
	status string,
) *GroupAssignmentChapterItemBuilder {
	builder.Status = &status
	return builder
}

func (builder *GroupAssignmentChapterItemBuilder) Build() models.GroupAssignmentChapterItem {
	return models.GroupAssignmentChapterItem{
		GroupID:        builder.GroupID,
		ChapterID:      builder.ChapterID,
		ItemID:         builder.ItemID,
		ExerciseIDList: builder.ExerciseIDList,
		FullMark:       builder.FullMark,
		TimeStart:      builder.TimeStart,
		TimeEnd:        builder.TimeEnd,
		Status:         builder.Status,
	}
}
