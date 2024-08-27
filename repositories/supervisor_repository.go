package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type SupervisorRepositoryQ interface{}

type SupervisorRepository struct {
	DB *gorm.DB
}

func NewSupervisorRepository(db *gorm.DB) *SupervisorRepository {
	return &SupervisorRepository{DB: db}
}

func (supervisorRepository *SupervisorRepository) GetAllSupervisors(
	supervisors *[]models.Supervisor,
) {
	supervisorRepository.DB.Preload("User").Find(supervisors)
}
