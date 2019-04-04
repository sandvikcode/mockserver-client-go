package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpectations(t *testing.T) {

	// Define test table
	testCases := []struct {
		description  string
		expectation  *Expectation
		expectedJSON string
	}{
		{"Path and method should be matched and then 200 returned.", CreateExpectation(WhenRequestPath("/path"), WhenRequestMethod("GET"), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path",
				"method": "GET"
			},
			"httpResponse": {
				"statusCode": 200
			}
		}`},
		{"Path should be matched and then 200 returned.", CreateExpectation(WhenRequestPath("/path"), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"httpResponse": {
				"statusCode": 200
			}
		}`},
		{"Path and query string parameters should be matched and then 200 returned.", CreateExpectation(WhenRequestPath("/path"), WhenRequestQueryStringParameters(map[string][]string{"name": {"value"}}), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path",
				"queryStringParameters" : {
					"name" : [ "value" ]
				}
			},
			"httpResponse": {
				"statusCode": 200
			}
		}`},
		{"Path and headers should be matched and then 200 returned.", CreateExpectation(WhenRequestPath("/path"), WhenRequestHeaders(map[string][]string{"name": {"value"}}), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path",
				"headers": {
					"name": ["value"]
				}
			},
			"httpResponse": {
				"statusCode": 200
			}
		}`},
		{"Path and authorization header should be matched and then 200 returned.", CreateExpectation(WhenRequestPath("/path"), WhenRequestAuth("mytoken"), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path",
				"headers": {
					"authorization": ["Bearer mytoken"]
				}
			},
			"httpResponse": {
				"statusCode": 200
			}
		}`},
		{"Path should be matched and then a specific JSON body returned.", CreateExpectation(WhenRequestPath("/path"), ThenResponseJSON(`{"key" : "value"}`)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"httpResponse": {
				"headers": {
					"content-type": ["application/json"]
				},
				"body": {
					"type" : "STRING",
					"string" : "{\"key\" : \"value\"}"
				}
			}
		}`},
		{"Path should be matched and then a specific text body returned.", CreateExpectation(WhenRequestPath("/path"), ThenResponseText(`Random string`)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"httpResponse": {
				"headers": {
					"content-type": ["text/plain; charset=utf-16"]
				},
				"body": {
					"type" : "STRING",
					"string" : "Random string"
				}
			}
		}`},
		{"Path and method should be matched and then 200 returned after 5 seconds.", CreateExpectation(WhenRequestPath("/path"), ThenResponseDelay(5*time.Second), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"httpResponse": {
				"statusCode": 200,
				"delay": {
					"timeUnit": "MILLISECONDS",
					"value": 5000
				}
			}
		}`},
		{"Path and method should be matched and then 200 returned but only 3 times.", CreateExpectation(WhenRequestPath("/path"), WhenTimes(3), ThenResponseStatus(http.StatusOK)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"httpResponse": {
				"statusCode": 200
			},
			"times": {
				"remainingTimes" : 3,
  				"unlimited" : false
			}
		}`},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a test server so we can inspect the JSON body sent by the mock-server client
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err, "Body reader must not return an error.")

				bodyMap := make(map[string]interface{})
				err = json.Unmarshal(body, &bodyMap)
				require.NoError(t, err, "Body un-marshall must not return an error.")

				expectedMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(tc.expectedJSON), &expectedMap)
				require.NoError(t, err, "Body un-marshall must not return an error.")

				require.Equal(t, expectedMap, bodyMap)

			}))
			defer ts.Close()

			mockClient := &Client{
				BaseURL: ts.URL,
				T:       t,
			}
			mockClient.AddExpectation(tc.expectation)
		})
	}
}
