package builders

import (
    "github.com/google/uuid"
    "github.com/Project-IPCA/ipca-backend/models"
)

type LabClassInfoBuilder struct {
    ChapterID    uuid.UUID
    ChapterIndex int
    Name         string
    FullMark     int
    NoItems      int
}

func NewLabClassInfoBuilder() *LabClassInfoBuilder {
    return &LabClassInfoBuilder{}
}

func (builder *LabClassInfoBuilder) SetChapterID(chapterID uuid.UUID) *LabClassInfoBuilder {
    builder.ChapterID = chapterID
    return builder
}

func (builder *LabClassInfoBuilder) SetChapterIndex(chapterIndex int) *LabClassInfoBuilder {
    builder.ChapterIndex = chapterIndex
    return builder
}

func (builder *LabClassInfoBuilder) SetName(name string) *LabClassInfoBuilder {
    builder.Name = name
    return builder
}

func (builder *LabClassInfoBuilder) SetFullMark(fullMark int) *LabClassInfoBuilder {
    builder.FullMark = fullMark
    return builder
}

func (builder *LabClassInfoBuilder) SetNoItems(noItems int) *LabClassInfoBuilder {
    builder.NoItems = noItems
    return builder
}

func (builder *LabClassInfoBuilder) Build() models.LabClassInfo {
    labclassInfo := models.LabClassInfo{
        ChapterID:    builder.ChapterID,
        ChapterIndex: builder.ChapterIndex,
        Name:         builder.Name,
        FullMark:     builder.FullMark,
        NoItems:      builder.NoItems,
    }
	return labclassInfo
}