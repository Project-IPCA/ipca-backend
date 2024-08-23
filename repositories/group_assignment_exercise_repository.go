package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupAssignmentExerciseRepositoryQ interface {
	GetGroupAssignmnetExercisesByGroupID(
		groupAssignmentExercises *[]models.GroupAssignmentExercise,
		groupId uuid.UUID,
	)
}

type GroupAssignmentExerciseRepository struct {
	DB *gorm.DB
}

func NewGroupAssignmentExerciseRepository(db *gorm.DB) *GroupAssignmentExerciseRepository {
	return &GroupAssignmentExerciseRepository{DB: db}
}

func (groupAssignmentExerciseRepository *GroupAssignmentExerciseRepository) GetGroupAssignmnetExercisesByGroupID(
	groupAssignmentExercises *[]models.GroupAssignmentExercise,
	groupId uuid.UUID,
) {
	groupAssignmentExerciseRepository.DB.Where("group_id = ?", groupId).
		Find(groupAssignmentExercises)
}
