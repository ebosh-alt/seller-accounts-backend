package domain

import (
	"errors"
	"strconv"
)

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, errors.New("invalid value")
	}
	return parsed, nil
}
