package detections

import (
	"net/http"

	"github.com/cfsdes/nucke/pkg/requests"
	"github.com/cfsdes/nucke/pkg/plugins/utils"
)


// Matcher Check main function
func MatchCheck(m Matcher, resp *http.Response, resTime int, oobID string, rawReq string, payload string, parameter string, resultChan chan Result) {
	// Get URL from raw request
	url := requests.ExtractRawURL(rawReq)
	
	foundArray := make([]bool, 0)
	statusCode, resBody, resHeaders, rawResp := requests.ParseResponse(resp)

	// Verifying if all rules matched
	if m.Body != nil {
		found := MatchBody(m.Body.RegexList, resBody)
		foundArray = append(foundArray, found)
	}
	if m.Header != nil {
		found := MatchHeader(m.Header.RegexList, resHeaders)
		foundArray = append(foundArray, found)
	}
	if m.ContentLength != nil {
		found := MatchMathOperation(m.ContentLength.Length, m.ContentLength.Operator, len(resBody))
		foundArray = append(foundArray, found)
	}
	if m.Time != nil {
		found := MatchMathOperation(m.Time.Seconds, m.Time.Operator, resTime)
		foundArray = append(foundArray, found)
	}
	if m.StatusCode != nil {
		found := MatchMathOperation(m.StatusCode.Code,m.StatusCode.Operator, statusCode)
		foundArray = append(foundArray, found)
	}
	if m.OOB {
		found := utils.CheckOobInteraction(oobID)
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
			resultChan <- Result{allTrue, url, payload, parameter, rawReq, rawResp, resBody}

		// OR condition
		} else if m.Operator == "OR" {
			for _, value := range foundArray {
				if value {
					resultChan <- Result{true, url, payload, parameter, rawReq, rawResp, resBody}
				}
			}
			resultChan <- Result{false, url, payload, parameter, rawReq, rawResp, resBody}
		}
	} else {
		resultChan <- Result{false, url, payload, parameter, rawReq, rawResp, resBody}
	}
}


