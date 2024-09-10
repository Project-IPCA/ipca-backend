package utils

import (
	"github.com/google/uuid"
)

func ContainsUUID(ids []uuid.UUID, id uuid.UUID) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
