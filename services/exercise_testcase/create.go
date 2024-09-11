package exercisetestcase

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func(exerciseTestcaseService *Service) Create(
    exerciseID       uuid.UUID,
    isReady          string,
    testcaseContent  string,
    isActive         *bool,
    isShowStudent    *bool,
    testcaseNote     *string,
    testcaseOutput   *string,
    testcaseError    *string,
)(*models.ExerciseTestcase,error){
	testcaseID:= uuid.New()
    
	exerciseTestcase := builders.NewExerciseTestcaseBuilder().
        SetExerciseID(exerciseID).
		SetTestcaseID(testcaseID).
		SetIsReady(isReady).
        SetTestcaseContent(testcaseContent).
        SetIsActive(isActive).
        SetIsShowStudent(isShowStudent).
        SetTestcaseNote(&testcaseContent).
        SetTestcaseOutput(testcaseOutput).
        SetTestcaseError(testcaseError).
        Build()
    
    exerciseTestcaseService.DB.Create(&exerciseTestcase)
    return exerciseTestcase,nil
}