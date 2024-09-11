package builders

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type ExerciseTestcaseBuilder struct {
    TestcaseID       uuid.UUID
    ExerciseID       uuid.UUID
    IsReady          string
    TestcaseContent  string
    IsActive         *bool
    IsShowStudent    *bool
    TestcaseNote     *string
    TestcaseOutput   *string
    TestcaseError    *string
}

func NewExerciseTestcaseBuilder() *ExerciseTestcaseBuilder {
    return &ExerciseTestcaseBuilder{}
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetTestcaseID(testcaseID uuid.UUID) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.TestcaseID = testcaseID
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetExerciseID(exerciseID uuid.UUID) *ExerciseTestcaseBuilder {
    fmt.Println("exerciseId")
    fmt.Println(exerciseID)
    exerciseTestcaseBuilder.ExerciseID = exerciseID
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetIsReady(isReady string) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.IsReady = isReady
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetTestcaseContent(content string) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.TestcaseContent = content
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetIsActive(isActive *bool) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.IsActive = isActive
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetIsShowStudent(isShowStudent *bool) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.IsShowStudent = isShowStudent
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetTestcaseNote(note *string) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.TestcaseNote = note
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetTestcaseOutput(output *string) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.TestcaseOutput = output
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) SetTestcaseError(err *string) *ExerciseTestcaseBuilder {
    exerciseTestcaseBuilder.TestcaseError = err
    return exerciseTestcaseBuilder
}

func (exerciseTestcaseBuilder *ExerciseTestcaseBuilder) Build() *models.ExerciseTestcase {
    exerciseTestcase := models.ExerciseTestcase{
        TestcaseID:       &exerciseTestcaseBuilder.TestcaseID,
        ExerciseID:       exerciseTestcaseBuilder.ExerciseID,
        IsReady:          exerciseTestcaseBuilder.IsReady,
        TestcaseContent:  exerciseTestcaseBuilder.TestcaseContent,
        IsActive:         exerciseTestcaseBuilder.IsActive,
        IsShowStudent:    exerciseTestcaseBuilder.IsShowStudent,
        TestcaseNote:     exerciseTestcaseBuilder.TestcaseNote,
        TestcaseOutput:   exerciseTestcaseBuilder.TestcaseOutput,
        TestcaseError:    exerciseTestcaseBuilder.TestcaseError,
    }
	return &exerciseTestcase
}