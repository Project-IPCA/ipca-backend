package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupAssignmentExerciseBuilder struct {
	GroupID    uuid.UUID
	ExerciseID uuid.UUID
	Selected   bool
}

func NewGroupAssignmentExerciseBuilder() *GroupAssignmentExerciseBuilder {
	return &GroupAssignmentExerciseBuilder{}
}

func (gae *GroupAssignmentExerciseBuilder) SetGroupID(
	groupId uuid.UUID,
) *GroupAssignmentExerciseBuilder {
	gae.GroupID = groupId
	return gae
}

func (gae *GroupAssignmentExerciseBuilder) SetExerciseID(
	exerciseId uuid.UUID,
) *GroupAssignmentExerciseBuilder {
	gae.ExerciseID = exerciseId
	return gae
}

func (gae *GroupAssignmentExerciseBuilder) SetSelected(
	selected bool,
) *GroupAssignmentExerciseBuilder {
	gae.Selected = selected
	return gae
}

func (gae *GroupAssignmentExerciseBuilder) Build() models.GroupAssignmentExercise {
	return models.GroupAssignmentExercise{
		GroupID:    gae.GroupID,
		ExerciseID: gae.ExerciseID,
		Selected:   gae.Selected,
	}
}
