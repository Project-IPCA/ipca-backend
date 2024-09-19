package studentassignmentchapteritem

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

func (studentAssignmentItemService *Service) UpdateAssignExercise (
	studentAssignmentChapterItem *models.StudentAssignmentChapterItem,
	exerciseId *uuid.UUID,
	){
		studentAssignmentItemService.DB.Model(studentAssignmentChapterItem).Update("exercise_id",exerciseId)
		studentAssignmentChapterItem.ExerciseID = exerciseId
}