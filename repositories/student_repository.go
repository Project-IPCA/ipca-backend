package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type StudentRepositoryQ interface{}

type StudentRepository struct {
	DB *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (studentRepository *StudentRepository) GetUserByKmitlID(
	student *models.Student,
	kmitlId string,
) {
	studentRepository.DB.Where("kmitl_id = ?", kmitlId).Find(student)
}
