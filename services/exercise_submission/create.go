package exercisesubmission

import (
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (exerciseSubmissionService *Service) Create(
	stuID uuid.UUID,
	exerciseID uuid.UUID,
	sourcecodeFilename string,
	marking int,
	isInfLoop *bool,
	output *string,
	result *string,
	errorMessage *string,
) (uuid.UUID, error) {
	submissionId := uuid.New()
	submission := builders.NewExerciseSubmissionBuilder().
	SetSubmissionID(submissionId).
	SetStuID(stuID).
	SetExerciseID(exerciseID).
	SetStatus("PENDING").
	SetSourcecodeFilename(sourcecodeFilename).
	SetMarking(marking).
	SetIsInfLoop(isInfLoop).
	SetOutput(output).
	SetResult(result).
	SetErrorMessage(errorMessage).
	Build()

	exerciseSubmissionService.DB.Create(&submission)
	return submissionId, nil
}