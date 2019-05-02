package mockclient

// VerificationSequence defines a specific ordered sequence of requests to MockServer
type VerificationSequence struct {
	Requests []*RequestMatcher `json:"httpRequests"`
}

// CreateVerification converts a number of expectation parts (options) into a single Expectation
func CreateVerification(opts ...ExpectationOption) *Expectation {
	// Specify some defaults if no options are set
	e := &Expectation{
		Request: &RequestMatcher{
			Path: "/(.*)",
		},
		Times: &Times{
			AtLeast: integerPointer(1),
			AtMost:  integerPointer(1),
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
		v.Times.AtLeast = integerPointer(times)
		return v
	}
}

// ThenAtMostCalls creates a verification that a matching call was received at most x times by MockServer
func ThenAtMostCalls(times int) ExpectationOption {
	return func(v *Expectation) *Expectation {
		v.Times.AtMost = integerPointer(times)
		return v
	}
}

// CreateVerificationSequence creates verifications for a given expectation sequence
func CreateVerificationSequence(opts ...ExpectationOption) *VerificationSequence {

	eArray := make([]*Expectation, 0)
	for _, opt := range opts {
		e := &Expectation{
			Request: &RequestMatcher{},
		}
		eArray = append(eArray, opt(e))
	}
	//TODO: tidy up
	vs := &VerificationSequence{}
	for _, item := range eArray {
		vs.Requests = append(vs.Requests, item.Request)
	}
	return vs
}

func integerPointer(i int) *int {
	return &i
}
