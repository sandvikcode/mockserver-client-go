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
	return func(e *Expectation) *Expectation {
		e.Times.AtLeast = integerPointer(times)
		return e
	}
}

// ThenAtMostCalls creates a verification that a matching call was received at most x times by MockServer
func ThenAtMostCalls(times int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Times.AtMost = integerPointer(times)
		return e
	}
}

func integerPointer(i int) *int {
	return &i
}
