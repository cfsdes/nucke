package fuzzers

import (
	"net/http"

    "github.com/cfsdes/nucke/pkg/plugins/detections"
)

func FuzzAll(r *http.Request, client *http.Client, payloads []string, matcher detections.Matcher) (bool, string, string, string, string, string, []detections.Result) {

	var logScansCombined []detections.Result

	allfuzz := []func(*http.Request, *http.Client, []string, detections.Matcher) (bool, string, string, string, string, string, []detections.Result){
		FuzzJSON,
		FuzzQuery,
		FuzzFormData,
		FuzzXML,
	}
	
	for _, fuzzer := range allfuzz {
		found, url, payload, param, rawReq, rawResp, logsScan := fuzzer(r, client, payloads, matcher)
		logScansCombined = append(logScansCombined, logsScan...)
		if found {
			return found, url, payload, param, rawReq, rawResp, logScansCombined
		}
	}

	return false, "", "", "", "", "", logScansCombined
}
