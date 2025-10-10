package validation

import (
	"errors"
	"regexp"
)

// Username validation rules
const (
	MinUsernameLength = 3
	MaxUsernameLength = 16
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]*[a-zA-Z][a-zA-Z0-9_]*$`)

// check that the username length is between Min/Max + contains only characters, numbers and underscores. (at least one letter)
func ValidateUsername(username string) error {
	if (len(username) < MinUsernameLength || len(username) > MaxUsernameLength) || !usernameRegex.MatchString(username) {
		return errors.New("the username must be between 3 and 16 characters and contain only alphanumeric characters or underscores, at least one letter")
	}

	return nil
}
