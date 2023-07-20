package detections

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
	RawResp string
	ResBody string
}