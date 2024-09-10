package groupassignmentexercise

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (groupAssignmentExerciseService *Service) Create(
	groupId uuid.UUID,
	exerciseId uuid.UUID,
	selected bool,
) error {
	gae := builders.NewGroupAssignmentExerciseBuilder().
		SetGroupID(groupId).
		SetExerciseID(exerciseId).
		SetSelected(selected).
		Build()

	groupAssignmentExerciseService.DB.Create(&gae)
	return nil
}
