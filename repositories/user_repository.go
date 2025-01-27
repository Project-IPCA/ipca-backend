package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
)

type UserRepositoryQ interface{}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepository *UserRepository) GetUserByUsername(
	user *models.User,
	username string,
) {
	userRepository.DB.Preload("Student").Where("username = ?", username).Find(user)
}

func (userRepository *UserRepository) GetUserByUserID(
	user *models.User,
	userId uuid.UUID,
) {
	userRepository.DB.Preload("Student").
		Preload("Supervisor").
		Preload("Dept").
		Where("user_id = ?", userId).
		Find(user)
}

func (userRepository *UserRepository) GetUserStudentAndGroupByUserID(
	user *models.User,
	userId uuid.UUID,
) {
	userRepository.DB.Preload("Student.Group").
		Preload("Supervisor").
		Preload("Dept").
		Where("user_id = ?", userId).
		Find(user)
}

func (userRepository *UserRepository) GetUserAdminRole(
	user *[]models.User,
) {
	userRepository.DB.Preload("Executive").
		Preload("Supervisor").
		Preload("TA").
		Where("role = ? OR role = ? OR role = ?", constants.Role.Supervisor, constants.Role.Ta, constants.Role.Executive).
		Find(user)
}
