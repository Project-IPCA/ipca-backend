package ta

import (
	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (taService *Service)CreateTa(
	taId uuid.UUID,
	groupId *uuid.UUID,
	note *string,
) error {
	ta := builders.NewTaBuilder().
		SetTaID(taId).
		SetGroupID(groupId).
		SetNote(note).
		Build()

	taService.DB.Create(ta)
	return nil
}