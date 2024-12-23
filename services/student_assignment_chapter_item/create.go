package studentassignmentchapteritem

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (studentAssignmentChapterItemService *Service) Create(
	stuID uuid.UUID,
	chapterID uuid.UUID,
	itemID int,
	exerciseID *uuid.UUID,
	fullMark int,
	marking int,
	timeStart *string,
	timeEnd *string,
) (*models.StudentAssignmentChapterItem, error) {

	studentAssignmentChapterItem := builders.NewStudentAssignmentChapterItemBuilder().
		SetStuID(stuID).
		SetChapterID(chapterID).
		SetItemID(itemID).
		SetExerciseID(exerciseID).
		SetFullMark(fullMark).
		SetMarking(marking).
		SetTimeStart(timeStart).
		SetTimeEnd(timeEnd).
		Build()
	err := studentAssignmentChapterItemService.DB.Create(&studentAssignmentChapterItem)
	if err.Error != nil {
		return nil, err.Error
	}
	return &studentAssignmentChapterItem, nil
}
