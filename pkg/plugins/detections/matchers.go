package detections

import (
	"regexp"
	"fmt"
)

// Check if regexList match with response body
func MatchBody(regexList []string, body string) bool {
	// Check if match some regex in the list (case insensitive)
	for _, regex := range regexList {
		match, err := regexp.MatchString("(?i)"+regex, string(body))
		if err != nil {
			fmt.Println("Could not check the match with regex: %v", err)
			return false
		}
		if match {
			return true
		}
	}

	return false
}

// Check if regexList match with response headers
func MatchHeader(regexList []string, responseHeaders map[string][]string) bool {
	// Check if response headers match with regexList
    for _, regex := range regexList {
        for k, v := range responseHeaders {
            if matched, _ := regexp.MatchString(regex, k+": "+v[0]); matched {
                return true
                break
            }
        }
    }

	return false
}

// Check if number match with operator
func MatchMathOperation(value int, operator string, target int) bool {
    switch operator {
    case "<":
        return target < value
    case ">":
        return target > value
    case "==":
        return target == value
	case "!=":
		return target != value
	case "<=":
		return target <= value
	case ">=":
		return target >= value
    default:
        return false
    }
}