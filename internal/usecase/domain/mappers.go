package domain

import (
	"strconv"
)

func parseCategoryID(rawID string) (int, error) {
	if rawID == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(rawID)
	if err != nil || parsed < 0 {
		return 0, ErrInvalidAccount
	}
	return parsed, nil
}
