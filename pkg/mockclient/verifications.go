package mockclient

// VerificationSequence defines a specific sequence of calls to MockServer
type VerificationSequence struct {
	Path string `json:"path,omitempty"`
}

// CreateVerification converts a number of expectation parts (options) into a single Expectation
func CreateVerification(opts ...ExpectationOption) *Expectation {
	// Specify some defaults if no options are set
	e := &Expectation{
		Request: &RequestMatcher{
			Path: "/(.*)",
		},
		Times: &Times{
			AtLeast: 1,
			AtMost:  1,
		},
	}
	// Append all options that are set (discard defaults)
	for _, opt := range opts {
		e = opt(e)
	}

	return e
}

// ThenAtLeastCalls creates a verification that a matching call was received at least x times by MockServer
func ThenAtLeastCalls(times int) ExpectationOption {
	return func(v *Expectation) *Expectation {
		v.Times.AtLeast = times
		return v
	}
}

// ThenAtMostCalls creates a verification that a matching call was received at most x times by MockServer
func ThenAtMostCalls(times int) ExpectationOption {
	return func(v *Expectation) *Expectation {
		v.Times.AtMost = times
		return v
	}
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
