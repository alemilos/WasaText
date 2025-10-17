package validation

import (
	"errors"
)

const (
	MinGroupLength = 1
	MaxGroupLength = 30
)

func ValidateGroupName(groupName string) error {
	if len(groupName) < MinGroupLength || len(groupName) > MaxGroupLength {
		return errors.New("the group name must be between 1 and 30 characters long")
	}

	return nil
}
