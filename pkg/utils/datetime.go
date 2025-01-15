package utils

import "time"

func IsTimeInRange(startTime, endTime *time.Time) bool {
	now := time.Now()

	if startTime == nil && endTime == nil {
		return true
	}

	if startTime != nil && endTime == nil {
		return now.After(*startTime) || now.Equal(*startTime)
	}

	if startTime == nil && endTime != nil {
		return now.Before(*endTime) || now.Equal(*endTime)
	}

	return (now.After(*startTime) || now.Equal(*startTime)) &&
		(now.Before(*endTime) || now.Equal(*endTime))
}
