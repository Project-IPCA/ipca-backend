package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
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
	userRepository.DB.Where("username = ?", username).Find(user)
}
