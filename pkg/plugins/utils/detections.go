package utils

import (
	"regexp"
	"io/ioutil"
	"fmt"
	"net/http"

	internalUtils "github.com/cfsdes/nucke/internal/utils"
)

// Matcher Structure
type Matcher struct {
    Time   *TimeMatcher
	StatusCode *StatusCodeMatcher
    Body   *BodyMatcher
    Header *HeaderMatcher
	OOB    bool
}

type TimeMatcher struct {
    Operator string
    Seconds  int
}

type StatusCodeMatcher struct {
    Operator string
    Code  int
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
}

// Matcher Check main function
func MatchChek(m Matcher, resp *http.Response, resTime int, oobID string, rawReq string, resultChan chan Result) {
	// Get URL from raw request
	url := ExtractRawURL(rawReq)

	foundArray := make([]bool, 0)
	statusCode, resBody, resHeaders := parseResponse(resp)

	// Verifying if all rules matched
	if m.Body != nil {
		found := matchBody(m.Body.RegexList, resBody)
		foundArray = append(foundArray, found)
	}
	if m.Header != nil {
		found := matchHeader(m.Header.RegexList, resHeaders)
		foundArray = append(foundArray, found)
	}
	if m.Time != nil {
		found := matchTime(m.Time.Seconds, m.Time.Operator, resTime)
		foundArray = append(foundArray, found)
	}
	if m.StatusCode != nil {
		found := matchStatusCode(m.StatusCode.Code,m.StatusCode.Operator, statusCode)
		foundArray = append(foundArray, found)
	}
	if m.OOB && oobID != "" {
		found := internalUtils.CheckOobInteraction(oobID)
		foundArray = append(foundArray, found)
	}

	// Validate if all matches are true
    if len(foundArray) > 0 {
		allTrue := true
		for _, value := range foundArray {
			if value == false {
				allTrue = false
				break
			}
		}
		resultChan <- Result{allTrue, rawReq, url}
	}

	resultChan <- Result{false, "", ""}
}

// Parse response
func parseResponse(resp *http.Response) (int, string, map[string][]string) {
	// Get status code
    statusCode := resp.StatusCode

	// Get response body
	responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Match parser error: ", err)
    }
	defer resp.Body.Close()

	// Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

	return statusCode, string(responseBody), responseHeaders
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

// Check if time match with time rule
func matchTime(seconds int, operator string, resTime int) bool {
    switch operator {
    case "<":
        return resTime < seconds
    case ">":
        return resTime > seconds
    case "==":
        return resTime == seconds
    default:
        return false
    }
}

// Check if status code match
func matchStatusCode(statusCode int, operator string, respStatusCode int) bool {
    switch operator {
    case "<":
        return respStatusCode < statusCode
    case ">":
        return respStatusCode > statusCode
    case "==":
        return respStatusCode == statusCode
    default:
        return false
    }
}