package mock

// VerifyTimes defines how many times the Mock-Server can be called
type VerifyTimes struct {
	AtLeast int `json:"atLeast,omitempty"`
	AtMost  int `json:"atMost,omitempty"`
}

// Verification defines how many times the Mock-Server can be called for a given request pattern
type Verification struct {
	Request *RequestMatcher `json:"httpRequest"`
	Times   *VerifyTimes    `json:"times,omitempty"`
}

// CreateVerification creates a verification for a given expectation
func CreateVerification(e *Expectation, vt *VerifyTimes) *Verification {
	v := &Verification{
		Request: e.Request,
		Times:   vt,
	}
	return v
}

// CreateVerify creates a verification that bounds the number of times MockServer should have been called
// Note: -1 is infinite number of times.
func CreateVerify(min, max int) *VerifyTimes {
	vt := &VerifyTimes{
		AtLeast: min,
		AtMost:  max,
	}
	return vt
}
