package builders

import (
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type ExerciseSubmissionBuilder struct {
    SubmissionID        uuid.UUID
    StuID               uuid.UUID
    ExerciseID          uuid.UUID
    Status              string
    SourcecodeFilename  string
    Marking             int
    TimeSubmit          time.Time
    IsInfLoop           *bool
    Output              *string
    Result              *string
    ErrorMessage        *string
}

func NewExerciseSubmissionBuilder() *ExerciseSubmissionBuilder {
    return &ExerciseSubmissionBuilder{}
}

func (b *ExerciseSubmissionBuilder) SetSubmissionID(submissionID uuid.UUID) *ExerciseSubmissionBuilder {
    b.SubmissionID = submissionID
    return b
}

func (b *ExerciseSubmissionBuilder) SetStuID(stuID uuid.UUID) *ExerciseSubmissionBuilder {
    b.StuID = stuID
    return b
}

func (b *ExerciseSubmissionBuilder) SetExerciseID(exerciseID uuid.UUID) *ExerciseSubmissionBuilder {
    b.ExerciseID = exerciseID
    return b
}

func (b *ExerciseSubmissionBuilder) SetStatus(status string) *ExerciseSubmissionBuilder {
    b.Status = status
    return b
}

func (b *ExerciseSubmissionBuilder) SetSourcecodeFilename(filename string) *ExerciseSubmissionBuilder {
    b.SourcecodeFilename = filename
    return b
}

func (b *ExerciseSubmissionBuilder) SetMarking(marking int) *ExerciseSubmissionBuilder {
    b.Marking = marking
    return b
}

func (b *ExerciseSubmissionBuilder) SetTimeSubmit(timeSubmit time.Time) *ExerciseSubmissionBuilder {
    b.TimeSubmit = timeSubmit
    return b
}

func (b *ExerciseSubmissionBuilder) SetIsInfLoop(isInfLoop *bool) *ExerciseSubmissionBuilder {
    b.IsInfLoop = isInfLoop
    return b
}

func (b *ExerciseSubmissionBuilder) SetOutput(output *string) *ExerciseSubmissionBuilder {
    b.Output = output
    return b
}

func (b *ExerciseSubmissionBuilder) SetResult(result *string) *ExerciseSubmissionBuilder {
    b.Result = result
    return b
}

func (b *ExerciseSubmissionBuilder) SetErrorMessage(errorMessage *string) *ExerciseSubmissionBuilder {
    b.ErrorMessage = errorMessage
    return b
}

func (b *ExerciseSubmissionBuilder) Build() models.ExerciseSubmission {
    return models.ExerciseSubmission{
        SubmissionID:       b.SubmissionID,
        StuID:              b.StuID,
        ExerciseID:         b.ExerciseID,
        Status:             b.Status,
        SourcecodeFilename: b.SourcecodeFilename,
        Marking:            b.Marking,
        TimeSubmit:         b.TimeSubmit,
        IsInfLoop:          b.IsInfLoop,
        Output:             b.Output,
        Result:             b.Result,
        ErrorMessage:       b.ErrorMessage,
    }
}