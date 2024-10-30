package detections

import (
	"net/http"

	"github.com/cfsdes/nucke/pkg/plugins/utils"
	"github.com/cfsdes/nucke/pkg/requests"
)

// Matcher Check main function
func MatchCheck(pluginDir string, m Matcher, resp *http.Response, resTime int, oobID string, rawReq string, payload string, parameter string) Result {
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
		found := MatchMathOperation(m.StatusCode.Code, m.StatusCode.Operator, statusCode)
		foundArray = append(foundArray, found)
	}
	if m.OOB {
		oob_id := utils.ExtractOobID(payload)
		utils.StoreDetection(pluginDir, oob_id, url, payload, parameter, rawReq)
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
			return Result{allTrue, url, payload, parameter, rawReq, rawResp, resBody}

			// OR condition
		} else if m.Operator == "OR" {
			for _, value := range foundArray {
				if value {
					return Result{true, url, payload, parameter, rawReq, rawResp, resBody}
				}
			}
			return Result{false, url, payload, parameter, rawReq, rawResp, resBody}
		}
	}

	return Result{false, url, payload, parameter, rawReq, rawResp, resBody}
}
