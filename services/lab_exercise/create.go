package labexercise

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (labExerciseService *Service) Create(
	request *requests.CreateLabExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetSourcecode(&request.Sourcecode).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}

func (labExerciseService *Service) CreateWithoutSourceCode(
	request *requests.CreateLabExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}
