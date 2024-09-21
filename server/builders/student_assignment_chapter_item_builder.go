package builders

import (
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type StudentAssignmentChapterItemsBuilder struct {
    StuID      uuid.UUID
    ChapterID  uuid.UUID
    ItemID     int
    ExerciseID *uuid.UUID
    FullMark   int
    Marking    int
    AddedDate  time.Time
    TimeStart  *string
    TimeEnd    *string
}

func NewStudentAssignmentChapterItemBuilder() *StudentAssignmentChapterItemsBuilder {
    return &StudentAssignmentChapterItemsBuilder{}
}

func (b *StudentAssignmentChapterItemsBuilder) SetStuID(stuID uuid.UUID) *StudentAssignmentChapterItemsBuilder {
    b.StuID = stuID
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetChapterID(chapterID uuid.UUID) *StudentAssignmentChapterItemsBuilder {
    b.ChapterID = chapterID
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetItemID(itemID int) *StudentAssignmentChapterItemsBuilder {
    b.ItemID = itemID
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetExerciseID(exerciseID *uuid.UUID) *StudentAssignmentChapterItemsBuilder {
    b.ExerciseID = exerciseID
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetFullMark(fullMark int) *StudentAssignmentChapterItemsBuilder {
    b.FullMark = fullMark
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetMarking(marking int) *StudentAssignmentChapterItemsBuilder {
    b.Marking = marking
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetAddedDate(addedDate time.Time) *StudentAssignmentChapterItemsBuilder {
    b.AddedDate = addedDate
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetTimeStart(timeStart *string) *StudentAssignmentChapterItemsBuilder {
    b.TimeStart = timeStart
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) SetTimeEnd(timeEnd *string) *StudentAssignmentChapterItemsBuilder {
    b.TimeEnd = timeEnd
    return b
}

func (b *StudentAssignmentChapterItemsBuilder) Build() models.StudentAssignmentChapterItem {
    return models.StudentAssignmentChapterItem{
        StuID:      b.StuID,
        ChapterID:  b.ChapterID,
        ItemID:     b.ItemID,
        ExerciseID: b.ExerciseID,
        FullMark:   b.FullMark,
        Marking:    b.Marking,
        AddedDate:  b.AddedDate,
        TimeStart:  b.TimeStart,
        TimeEnd:    b.TimeEnd,
    }
}