package mock

import (
	"fmt"
)

// Expectation defines the complete request/response interaction for a given scenario
type Expectation struct {
	Request  *RequestMatcher `json:"httpRequest"`
	Response *ActionResponse `json:"httpResponse,omitempty"`
	Times    *Times          `json:"times,omitempty"`
}

// RequestMatcher is used to match which requests the expectation will be applied to
type RequestMatcher struct {
	Path                  string              `json:"path,omitempty"`
	Method                string              `json:"method,omitempty"`
	Headers               map[string][]string `json:"headers,omitempty"`
	QueryStringParameters map[string][]string `json:"queryStringParameters,omitempty"`
}

// ActionResponse defines what actions to take when a request is matched e.g. response, delay, forward etc.
type ActionResponse struct {
	Headers    map[string][]string `json:"headers,omitempty"`
	StatusCode int                 `json:"statusCode,omitempty"`
	Body       *ResponseBody       `json:"body,omitempty"`
	Delay      *Delay              `json:"delay,omitempty"`
}

// ResponseBody sets the request body the MockServer will return when serving a matched response
type ResponseBody struct {
	Type   string `json:"type"`
	String string `json:"string"`
}

// Times sets how many times the MockServer will serve a given request
type Times struct {
	RemainingTime int  `json:"remainingTimes"`
	Unlimited     bool `json:"unlimited"`
}

// Delay sets how long the MockServer will wait before serving a matched response
type Delay struct {
	TimeUnit string `json:"timeUnit"`
	Value    int    `json:"value"`
}

// ExpectationOption enables building expectations in many parts
type ExpectationOption func(e *Expectation) *Expectation

// CreateExpectation converts a number of expectation parts (options) into a single Expectation
func CreateExpectation(opts ...ExpectationOption) *Expectation {
	e := &Expectation{
		Request: &RequestMatcher{
			Path: "/(.*)",
		},
		Response: &ActionResponse{},
	}

	for _, opt := range opts {
		e = opt(e)
	}

	return e
}

// WhenRequestHeaders creates an expectation based on required request headers
func WhenRequestHeaders(headers map[string][]string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		if e.Request.Headers == nil {
			e.Request.Headers = make(map[string][]string)
		}

		for h, v := range headers {
			e.Request.Headers[h] = v
		}

		return e
	}
}

// WhenRequestQueryStringParameters creates an expectation based on required query string parameters
func WhenRequestQueryStringParameters(qsp map[string][]string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		if e.Request.QueryStringParameters == nil {
			e.Request.QueryStringParameters = make(map[string][]string)
		}

		for q, v := range qsp {
			e.Request.QueryStringParameters[q] = v
		}

		return e
	}
}

// WhenRequestAuth creates an expectation based on a required Authorization request header
func WhenRequestAuth(authToken string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		if e.Request.Headers == nil {
			e.Request.Headers = make(map[string][]string)
		}

		e.Request.Headers["authorization"] = []string{fmt.Sprintf("Bearer %s", authToken)}

		return e
	}
}

// WhenRequestPath creates an expectation based on a path
func WhenRequestPath(path string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Request.Path = path
		return e
	}
}

// WhenRequestMethod creates an expectation based on an HTTP method
func WhenRequestMethod(method string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Request.Method = method
		return e
	}
}

// WhenTimes creates an expectation bounded by a limited number of calls
func WhenTimes(times int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Times = &Times{
			RemainingTime: times,
			Unlimited:     false,
		}
		return e
	}
}

// ThenResponseStatus creates an action that returns an HTTP status code when a request is matched
func ThenResponseStatus(statusCode int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Response.StatusCode = statusCode
		return e
	}
}

// ThenResponseJSON creates an action that returns an HTTP body when a request is matched
func ThenResponseJSON(body string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		r := e.Response
		r.Body = &ResponseBody{
			Type:   "STRING",
			String: body,
		}

		if r.Headers == nil {
			r.Headers = make(map[string][]string)
		}
		r.Headers["content-type"] = []string{"application/json"}

		return e
	}
}

// ThenResponseDelay creates an action that delays returning an HTTP response when a request is matched
func ThenResponseDelay(delay int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		r := e.Response
		r.Delay = &Delay{
			TimeUnit: "SECONDS",
			Value:    delay,
		}
		return e
	}
}
