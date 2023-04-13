package utils

import (
	"regexp"
	
	"github.com/pkg/errors"
)

func MatchString(regexList []string, body string) (bool, error) {
	// Check if match some regex in the list (case insensitive)
	for _, regex := range regexList {
		match, err := regexp.MatchString("(?i)"+regex, string(body))
		if err != nil {
			return false, errors.Errorf("Could not check the match with regex")
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}