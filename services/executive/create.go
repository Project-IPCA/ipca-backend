package executive

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (executiveService *Service) Create(executiveId uuid.UUID) error {
	executive := builders.NewExecutiveBuilder().
		SetExecutiveID(executiveId).
		Build()

	executiveService.DB.Create(&executive)
	return nil
}
