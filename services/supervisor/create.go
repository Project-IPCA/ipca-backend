package supervisor

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (supervisorService *Service) Create(supervisorId uuid.UUID, dept string) error {
	supervisor := builders.NewSupervisorBuilder().
		SetSupervisorID(supervisorId).
		SetDept(&dept).
		Build()

	supervisorService.DB.Create(&supervisor)
	return nil
}
