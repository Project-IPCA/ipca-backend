package builders

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type TaBuilder struct {
    TaID    uuid.UUID
    GroupID *uuid.UUID
    Note    *string
}

func NewTaBuilder() *TaBuilder {
    return &TaBuilder{}
}

func (b *TaBuilder) SetTaID(taID uuid.UUID) *TaBuilder {
    b.TaID = taID
    return b
}

func (b *TaBuilder) SetGroupID(groupID *uuid.UUID) *TaBuilder {
    b.GroupID = groupID
    return b
}

func (b *TaBuilder) SetNote(note *string) *TaBuilder {
    b.Note = note
    return b
}

func (b *TaBuilder) Build() models.TA {
    return models.TA{
		TaID: b.TaID,
		GroupID: b.GroupID,
		Note: b.Note,
    }
}
