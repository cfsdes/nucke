package utils

import (
	"regexp"
	"io/ioutil"
	"fmt"
	"net/http"
)

// Matcher Structure
type Matcher struct {
    Time   *TimeMatcher
    Body   *BodyMatcher
    Header *HeaderMatcher
}

type TimeMatcher struct {
    Operator string
    Seconds  int
}

type BodyMatcher struct {
    RegexList []string
}

type HeaderMatcher struct {
    RegexList []string
}


// Matcher Check main function
func MatchChek(m Matcher, resp *http.Response, resTime int) (bool) {

	foundArray := make([]bool, 0)
	resBody, resHeaders := parseResponse(resp)

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

	// Validate if all matches are true
	allTrue := true
    for _, value := range foundArray {
        if value == false {
            allTrue = false
            break
        }
    }

	return allTrue
}

// Parse response
func parseResponse(resp *http.Response) (string, map[string][]string) {
	// Get response body
	responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Match parser error: %v", err)
    }

	// Get response headers
    responseHeaders := make(map[string][]string)
    for k, v := range resp.Header {
        responseHeaders[k] = v
    }

	return string(responseBody), responseHeaders
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