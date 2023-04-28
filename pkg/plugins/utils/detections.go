package utils

import (
	"regexp"
	"fmt"
	"net/http"

	internalUtils "github.com/cfsdes/nucke/internal/utils"
	"github.com/cfsdes/nucke/pkg/requests"
)

// Matcher Structure
type Matcher struct {
    Time   *TimeMatcher
	StatusCode *StatusCodeMatcher
    Body   *BodyMatcher
    Header *HeaderMatcher
	ContentLength *ContentLengthMatcher
	OOB    bool
	Operator string
}

type TimeMatcher struct {
    Operator string
    Seconds  int
}

type StatusCodeMatcher struct {
    Operator string
    Code  int
}

type ContentLengthMatcher struct {
	Operator string
	Length int
}

type BodyMatcher struct {
    RegexList []string
}

type HeaderMatcher struct {
    RegexList []string
}

// Result structure
type Result struct {
	Found bool
	RawReq string
	URL string
	Payload string
	Param string
}

// Matcher Check main function
func MatchChek(m Matcher, resp *http.Response, resTime int, oobID string, rawReq string, payload string, parameter string, resultChan chan Result) {
	// Get URL from raw request
	url := requests.ExtractRawURL(rawReq)
	
	foundArray := make([]bool, 0)
	statusCode, resBody, resHeaders := requests.ParseResponse(resp)

	// Verifying if all rules matched
	if m.Body != nil {
		found := matchBody(m.Body.RegexList, resBody)
		foundArray = append(foundArray, found)
	}
	if m.Header != nil {
		found := matchHeader(m.Header.RegexList, resHeaders)
		foundArray = append(foundArray, found)
	}
	if m.ContentLength != nil {
		found := matchMathOperation(m.ContentLength.Length, m.ContentLength.Operator, len(resBody))
		foundArray = append(foundArray, found)
	}
	if m.Time != nil {
		found := matchMathOperation(m.Time.Seconds, m.Time.Operator, resTime)
		foundArray = append(foundArray, found)
	}
	if m.StatusCode != nil {
		found := matchMathOperation(m.StatusCode.Code,m.StatusCode.Operator, statusCode)
		foundArray = append(foundArray, found)
	}
	if m.OOB && oobID != "" {
		found := internalUtils.CheckOobInteraction(oobID)
		foundArray = append(foundArray, found)
	}

	// Validate if all matches are true
    if len(foundArray) > 0 {
		// AND condition
		if m.Operator == "" || m.Operator == "AND" {
			allTrue := true
			for _, value := range foundArray {
				if value == false {
					allTrue = false
					break
				}
			}
			resultChan <- Result{allTrue, rawReq, url, payload, parameter}

		// OR condition
		} else if m.Operator == "OR" {
			for _, value := range foundArray {
				if value {
					resultChan <- Result{true, rawReq, url, payload, parameter}
				}
			}
			resultChan <- Result{false, rawReq, url, payload, parameter}
		}
	} else {
		resultChan <- Result{false, "", "", "", ""}
	}
}

// Check if regexList match with response body
func matchBody(regexList []string, body string) bool {
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
func matchHeader(regexList []string, responseHeaders map[string][]string) bool {
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
func matchMathOperation(value int, operator string, target int) bool {
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
