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

func FindMax(arr []map[string]int, key string) int {
    if len(arr) == 0 {
        return 0
    }
    
    max := arr[0][key]
    for _, obj := range arr[1:] {
        if obj[key] > max {
            max = obj[key]
        }
    }
    return max
}
