package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ExecutiveBuilder struct {
	ExecutiveID uuid.UUID
}

func NewExecutiveBuilder() *ExecutiveBuilder {
	return &ExecutiveBuilder{}
}

func (executiveBuilder *ExecutiveBuilder) SetExecutiveID(
	executiveId uuid.UUID,
) (s *ExecutiveBuilder) {
	executiveBuilder.ExecutiveID = executiveId
	return executiveBuilder
}

func (executiveBuilder *ExecutiveBuilder) Build() models.Executive {
	executive := models.Executive{
		ExecutiveID: executiveBuilder.ExecutiveID,
	}
	return executive
}
