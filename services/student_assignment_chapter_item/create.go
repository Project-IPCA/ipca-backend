package studentassignmentchapteritem

import (
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (studentAssignmentChapterItemService *Service) Create (
	stuID      uuid.UUID,
    chapterID  uuid.UUID,
    itemID     uuid.UUID,
    exerciseID *uuid.UUID,
    fullMark   int,
    marking    int,
    timeStart  *string,
    timeEnd    *string,
) error {
	studentAssignmentChapterItem := builders.NewStudentAssignmentChapterItemBuilder().
		SetStuID(stuID).
		SetChapterID(chapterID).
		SetItemID(itemID).
		SetExerciseID(*exerciseID).
		SetFullMark(fullMark).
		SetMarking(marking).
		SetTimeStart(*timeStart).
		SetTimeEnd(*timeEnd).
		Build()
	studentAssignmentChapterItemService.DB.Create(&studentAssignmentChapterItem)
	return nil
}