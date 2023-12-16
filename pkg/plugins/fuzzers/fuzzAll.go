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
		match, rawReq, url, payload, param, rawResp, logsScan := fuzzer(r, client, payloads, matcher)
		logScansCombined = append(logScansCombined, logsScan...)
		if match {
			return match, rawReq, url, payload, param, rawResp, logScansCombined
		}
	}

	return false, "", "", "", "", "", logScansCombined
}
