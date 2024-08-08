package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type UserStudentRepositoryQ interface {
	GetUserByStuID(
		userStudent *models.UserStudent,
		stuId string,
	)
}

type UserStudentRepository struct {
	DB *gorm.DB
}

func NewUserStudentRepository(db *gorm.DB) *UserStudentRepository {
	return &UserStudentRepository{DB: db}
}

func (userStudentRepository *UserStudentRepository) GetUserByStuID(
	userStudent *models.UserStudent,
	stuId string,
) {
	userStudentRepository.DB.Where("stu_stu_id = ?", stuId).Find(userStudent)
}
