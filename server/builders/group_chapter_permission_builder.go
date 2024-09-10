package builders

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupChapterPermissionBuilder struct {
	ClassID         uuid.UUID
	ChapterID       uuid.UUID
	AllowAccessType string
	AccessTimeStart *time.Time
	AccessTimeEnd   *time.Time
	AllowSubmitType string
	SubmitTimeStart *time.Time
	SubmitTimeEnd   *time.Time
	AllowSubmit     bool
	Status          string
	AllowAccess     bool
	TimeStart       *string
	TimeEnd         *string
}

func NewGroupChapterPermissionBuilder() *GroupChapterPermissionBuilder {
	return &GroupChapterPermissionBuilder{}
}

func (builder *GroupChapterPermissionBuilder) SetClassID(
	classID uuid.UUID,
) *GroupChapterPermissionBuilder {
	builder.ClassID = classID
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetChapterID(
	chapterID uuid.UUID,
) *GroupChapterPermissionBuilder {
	builder.ChapterID = chapterID
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAllowAccessType(
	allowAccessType string,
) *GroupChapterPermissionBuilder {
	builder.AllowAccessType = allowAccessType
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAccessTimeStart(
	accessTimeStart time.Time,
) *GroupChapterPermissionBuilder {
	builder.AccessTimeStart = &accessTimeStart
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAccessTimeEnd(
	accessTimeEnd time.Time,
) *GroupChapterPermissionBuilder {
	builder.AccessTimeEnd = &accessTimeEnd
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAllowSubmitType(
	allowSubmitType string,
) *GroupChapterPermissionBuilder {
	builder.AllowSubmitType = allowSubmitType
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetSubmitTimeStart(
	submitTimeStart time.Time,
) *GroupChapterPermissionBuilder {
	builder.SubmitTimeStart = &submitTimeStart
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetSubmitTimeEnd(
	submitTimeEnd time.Time,
) *GroupChapterPermissionBuilder {
	builder.SubmitTimeEnd = &submitTimeEnd
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAllowSubmit(
	allowSubmit bool,
) *GroupChapterPermissionBuilder {
	builder.AllowSubmit = allowSubmit
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetStatus(
	status string,
) *GroupChapterPermissionBuilder {
	builder.Status = status
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetAllowAccess(
	allowAccess bool,
) *GroupChapterPermissionBuilder {
	builder.AllowAccess = allowAccess
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetTimeStart(
	timeStart string,
) *GroupChapterPermissionBuilder {
	builder.TimeStart = &timeStart
	return builder
}

func (builder *GroupChapterPermissionBuilder) SetTimeEnd(
	timeEnd string,
) *GroupChapterPermissionBuilder {
	builder.TimeEnd = &timeEnd
	return builder
}

func (builder *GroupChapterPermissionBuilder) Build() models.GroupChapterPermission {
	gcp := models.GroupChapterPermission{
		ClassID:         builder.ClassID,
		ChapterID:       builder.ChapterID,
		AllowAccessType: builder.AllowAccessType,
		AccessTimeStart: builder.AccessTimeStart,
		AccessTimeEnd:   builder.AccessTimeEnd,
		AllowSubmitType: builder.AllowSubmitType,
		SubmitTimeStart: builder.SubmitTimeStart,
		SubmitTimeEnd:   builder.SubmitTimeEnd,
		AllowSubmit:     builder.AllowSubmit,
		Status:          builder.Status,
		AllowAccess:     builder.AllowAccess,
		TimeStart:       builder.TimeStart,
		TimeEnd:         builder.TimeEnd,
	}
	return gcp
}
