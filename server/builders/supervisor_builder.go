package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type SupervisorBuilder struct {
	SupervisorID uuid.UUID
	Dept         *string
}

func NewSupervisorBuilder() *SupervisorBuilder {
	return &SupervisorBuilder{}
}

func (supervisorBuilder *SupervisorBuilder) SetSupervisorID(
	supervisorId uuid.UUID,
) (s *SupervisorBuilder) {
	supervisorBuilder.SupervisorID = supervisorId
	return supervisorBuilder
}

func (supervisorBuilder *SupervisorBuilder) SetDept(dept *string) (s *SupervisorBuilder) {
	supervisorBuilder.Dept = dept
	return supervisorBuilder
}

func (supervisorBuilder *SupervisorBuilder) Build() models.Supervisor {
	supervisor := models.Supervisor{
		SupervisorID: supervisorBuilder.SupervisorID,
		Dept:         supervisorBuilder.Dept,
	}
	return supervisor
}
