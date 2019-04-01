package mock

import (
	"fmt"
	"net/http"
)

// ResponseBody sets the request body the Mock-Server will return when serving a matched response
type ResponseBody struct {
	Type   string `json:"type"`
	String string `json:"string"`
}

// Delay sets how long the Mock-Server will wait before serving a matched response
type Delay struct {
	TimeUnit string `json:"timeUnit"`
	Value    int    `json:"value"`
}

// Times sets how many times the Mock-Server will serve a given request
type Times struct {
	RemainingTime int  `json:"remainingTimes"`
	Unlimited     bool `json:"unlimited"`
}

type ActionResponse struct {
	Headers    map[string][]string `json:"headers,omitempty"`
	StatusCode int                 `json:"statusCode"`
	Body       *ResponseBody       `json:"body,omitempty"`
	Delay      *Delay              `json:"delay,omitempty"`
}

type RequestMatcher struct {
	Path    string              `json:"path,omitempty"`
	Method  string              `json:"method,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type Expectation struct {
	Request  *RequestMatcher `json:"httpRequest"`
	Response *ActionResponse `json:"httpResponse,omitEmpty"`
	Times    *Times          `json:"times"`
}

type ExpectationOption func(e *Expectation) *Expectation

func CreateExpectation(opts ...ExpectationOption) *Expectation {
	e := &Expectation{
		Request: &RequestMatcher{
			Path: "/(.*)",
		},
		Response: &ActionResponse{
			StatusCode: http.StatusOK,
		},
	}

	for _, opt := range opts {
		e = opt(e)
	}

	return e
}

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

func WhenRequestAuth(authToken string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		if e.Request.Headers == nil {
			e.Request.Headers = make(map[string][]string)
		}

		e.Request.Headers["authorization"] = []string{fmt.Sprintf("Bearer %s", authToken)}

		return e
	}
}

func WhenRequestPath(path string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Request.Path = path
		return e
	}
}

func WhenRequestMethod(method string) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Request.Method = method
		return e
	}
}

func WhenTimes(times int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Times = &Times{
			RemainingTime: times,
			Unlimited:     false,
		}
		return e
	}
}

func ThenResponseStatus(statusCode int) ExpectationOption {
	return func(e *Expectation) *Expectation {
		e.Response.StatusCode = statusCode
		return e
	}
}

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
