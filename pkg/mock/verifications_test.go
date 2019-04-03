package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifications(t *testing.T) {

	// Define test table
	testCases := []struct {
		description  string
		verification *Verification
		expectedJSON string
	}{
		{"Verify the MockServer was called at least 5 times, and at most 10 times, for a given path.", CreateVerification(CreateExpectation(WhenRequestPath("/path")), CreateVerify(5, 10)), `
		{
			"httpRequest": {
				"path": "/path"
			},
			"times": {
				"atLeast": 5,
        		"atMost": 10
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

				require.Equal(t, bodyMap, expectedMap)

			}))
			defer ts.Close()

			mockClient := &Client{
				BaseURL: ts.URL,
				T:       t,
			}
			mockClient.AddVerification(tc.verification)
		})
	}
}
/*
func TestVerificationSequence(t *testing.T) {

	// Define test table
	testCases := []struct {
		description          string
		verificationSequence []*VerificationSequence
		expectedJSON         string
	}{
		{"Verify the MockServer was called with these specific calls in this specific order.", CreateVerificationSequence(VerifyPath("/some/path/one"), VerifyPath("/some/path/two"), VerifyPath("/some/path/three")), `
		[
			{
			"path": "/some/path/one"
			},
			{
			"path": "/some/path/two"
			},
			{
			"path": "/some/path/three"
			}
		]`},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a test server so we can inspect the JSON body sent by the mock-server client
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err, "Body reader must not return an error.")

				//bodyMap := make(map[string]interface{})
				//y := make([]interface{},0)
				err = json.Unmarshal(body, &y)
				require.NoError(t, err, "Body un-marshall must not return an error.")

				//expectedMap := make(map[string]interface{})
				x := make([]interface{},0)
				err = json.Unmarshal([]byte(tc.expectedJSON), &x)
				require.NoError(t, err, "Body un-marshall must not return an error.")

				require.Equal(t, y, x)

			}))
			defer ts.Close()

			mockClient := &Client{
				BaseURL: ts.URL,
				T:       t,
			}
			mockClient.AddVerificationSequence(tc.verificationSequence)
		})
	}
}
*/
