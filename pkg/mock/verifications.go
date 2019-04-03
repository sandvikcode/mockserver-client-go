package mock

// VerifyTimes defines how many times the MockServer may be called
type VerifyTimes struct {
	AtLeast int `json:"atLeast,omitempty"`
	AtMost  int `json:"atMost,omitempty"`
}

// Verification defines how many times the MockServer may be called for a given request pattern
type Verification struct {
	Request *RequestMatcher `json:"httpRequest"`
	Times   *VerifyTimes    `json:"times,omitempty"`
}

// VerificationSequence defines a specific sequence of calls to MockServer
type VerificationSequence struct {
	Path string `json:"path,omitempty"`
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
/*
// VerificationOption enables building verifications in many parts
type VerificationOption func(e *VerificationSequence) *VerificationSequence

// CreateVerificationSequence creates a verification for a given expectation sequence
func CreateVerificationSequence(opts ...VerificationOption) []*VerificationSequence {
	vsArray := make([]*VerificationSequence, 0)
	for _, opt := range opts {
		v := &VerificationSequence{}
		vsArray = append(vsArray, opt(v))
	}

	return vsArray
}

func VerifyPath(path string) VerificationOption {
	return func(vs *VerificationSequence) *VerificationSequence {
		vs.Path = path
		return vs
	}
}
*/
