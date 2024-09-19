package studentassignmentchapteritem

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

func (studentAssignmentItemService *Service) UpdateAssignExercise (
	studentAssignmentChapterItem *models.StudentAssignmentChapterItem,
	exerciseId *uuid.UUID,
	){
		studentAssignmentChapterItem.ExerciseID = exerciseId
		studentAssignmentItemService.DB.Save(studentAssignmentChapterItem)
}