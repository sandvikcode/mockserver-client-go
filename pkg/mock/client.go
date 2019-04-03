package mock

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"net/http"
	"net/url"
	"strings"

	"github.com/stretchr/testify/require"
)

// Client wraps the testing context used for all interactions with MockServer
type Client struct {
	// testing type T
	T *testing.T
	// mock server base address
	BaseURL string
}

// AddExpectation adds an expectation based on a request matcher to MockServer
func (hms *Client) AddExpectation(exp *Expectation) {
	msg, err := json.Marshal(exp)
	if err != nil {
		require.NoError(hms.T, err,
			"Failed to serialize mock server expectation.")
	}

	hms.callMock("expectation", string(msg))
}

// AddVerification adds a verification of requests to MockServer
func (hms *Client) AddVerification(v *Verification) {
	msg, err := json.Marshal(v)
	if err != nil {
		require.NoError(hms.T, err,
			"Failed to serialize mock server verification.")
	}

	hms.callMock("verify", string(msg))
}

// Clear everything that matches a given path in MockServer
func (hms *Client) Clear(path string) {
	mockReqBody := fmt.Sprintf(`
			{
				"path": "%s"
			}
			`, path)
	hms.callMock("clear", mockReqBody)
}

// Reset the entire MockServer, clearing all state
func (hms *Client) Reset() {
	hms.callMock("reset", "")
}

func (hms *Client) callMock(mockAPI, mockReqBody string) {
	mockURL := fmt.Sprintf("%s/%s", hms.BaseURL, mockAPI)
	// check url is valid
	if _, err := url.ParseRequestURI(mockURL); err != nil {
		require.NoError(hms.T, err,
			fmt.Sprintf("'%s' is not a valid mock server URL", mockURL))
	}

	hc := &http.Client{
		// Set timeout to 5s instead of default 30s
		Timeout: time.Duration(5 * time.Second),
	}
	reader := strings.NewReader(mockReqBody)

	mockReq, err := http.NewRequest("PUT", mockURL, reader)
	if err != nil {
		require.NoError(hms.T, err, "Failed to create request to mock server.")
	}
	mockRes, err := hc.Do(mockReq)
	if err != nil {
		require.NoError(hms.T, err, "Failed to send request to mock server.")
	}
	// MockServer verification returns 202 on success and 406 on failure
	if mockRes.StatusCode != http.StatusAccepted {
		require.NoError(hms.T, err,
			fmt.Sprintf("Mock server verification did not meet expectations and failed with status: %s", mockRes.Status))
	}
}
