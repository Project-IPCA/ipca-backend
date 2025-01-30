package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
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

func (supervisorRepository *SupervisorRepository) CheckValidSuperID(
	superId uuid.UUID,
) bool {
	var supervisor models.Supervisor
	err := supervisorRepository.DB.Where("supervisor_id = ?", superId).First(&supervisor)
	return err.Error == nil
}
