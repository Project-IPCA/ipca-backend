package repositories

import (
	"github.com/Project-IPCA/ipca-backend/models"
	exercisetestcase "github.com/Project-IPCA/ipca-backend/services/exercise_testcase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseTestcaseRepository struct{
	DB *gorm.DB
}

func NewExerciseTestcaseRepository(db *gorm.DB) *ExerciseTestcaseRepository {
	return &ExerciseTestcaseRepository{DB: db}
}

func (repo *ExerciseTestcaseRepository) DeleteExerciseTetscaseID (testcaseList []string){
	repo.DB.Where("testcase_id IN ?", testcaseList).Delete(&models.ExerciseTestcase{})
}

func (repo *ExerciseTestcaseRepository) UpsertExerciseTestcaseID (testcase models.ExerciseTestcase) *uuid.UUID{
	if(testcase.TestcaseID != nil){
		repo.DB.Model(&models.ExerciseTestcase{}).
			Where("testcase_id = ? AND exercise_id = ?", testcase.TestcaseID, testcase.ExerciseID).
			Updates(testcase)
		return nil
	}else{
		exerciseTestcaseService := exercisetestcase.NewExerciseTestcaseService(repo.DB)
		saveTestcase,_ := exerciseTestcaseService.Create(
			testcase.ExerciseID,
			testcase.IsReady,
			testcase.TestcaseContent,
			testcase.IsActive,
			testcase.IsShowStudent,
			testcase.TestcaseNote,
			testcase.TestcaseOutput,
			testcase.TestcaseError,
		)
		return saveTestcase.TestcaseID
	}
}

func (repo *ExerciseTestcaseRepository) GetTestcasesByExerciseID (exerciseId uuid.UUID,testcases *[]models.ExerciseTestcase){
	repo.DB.Where("exercise_id = ? ",exerciseId).Find(testcases)
}