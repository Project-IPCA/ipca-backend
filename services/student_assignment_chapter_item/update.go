package studentassignmentchapteritem

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

func (studentAssignmentItemService *Service) UpdateAssignExercise(
	studentAssignmentChapterItem *models.StudentAssignmentChapterItem,
	exerciseId *uuid.UUID,
) {
	studentAssignmentItemService.DB.Model(studentAssignmentChapterItem).Update("exercise_id", exerciseId)
	studentAssignmentChapterItem.ExerciseID = exerciseId
}

func (studentAssignmentItemService *Service) ResetMarking(
	studentId uuid.UUID,
	chapterId uuid.UUID,
	ItemId string,
) {
	studentAssignmentItemService.DB.Model(models.StudentAssignmentChapterItem{}).Where("stu_id = ? AND chapter_id = ? AND item_id = ?", studentId, chapterId, ItemId).Update("marking", 0)
}
