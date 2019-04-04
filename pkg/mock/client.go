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
func (c *Client) AddExpectation(exp *Expectation) {
	msg, err := json.Marshal(exp)
	if err != nil {
		require.NoError(c.T, err,
			"Failed to serialize mock server expectation.")
	}

	c.callMock("expectation", string(msg))
}

// AddVerification adds a verification of requests to MockServer
func (c *Client) AddVerification(exp *Expectation) {
	msg, err := json.Marshal(exp)
	if err != nil {
		require.NoError(c.T, err,
			"Failed to serialize mock server verification.")
	}

	c.callMock("verify", string(msg))
}

// AddVerificationSequence adds a verification of a specific sequence of requests to MockServer
func (c *Client) AddVerificationSequence(v []*VerificationSequence) {
	msg, err := json.Marshal(v)
	if err != nil {
		require.NoError(c.T, err,
			"Failed to serialize mock server verification sequence.")
	}

	c.callMock("verifySequence", string(msg))
}

// Clear everything that matches a given path in MockServer
func (c *Client) Clear(path string) {
	mockReqBody := fmt.Sprintf(`
			{
				"path": "%s"
			}
			`, path)
	c.callMock("clear", mockReqBody)
}

// Reset the entire MockServer, clearing all state
func (c *Client) Reset() {
	c.callMock("reset", "")
}

func (c *Client) callMock(mockAPI, mockReqBody string) {
	mockURL := fmt.Sprintf("%s/%s", c.BaseURL, mockAPI)
	// check url is valid
	if _, err := url.ParseRequestURI(mockURL); err != nil {
		require.NoError(c.T, err,
			fmt.Sprintf("'%s' is not a valid mock server URL", mockURL))
	}

	hc := &http.Client{
		// Set timeout to 5s instead of default 30s
		Timeout: time.Duration(5 * time.Second),
	}
	reader := strings.NewReader(mockReqBody)

	mockReq, err := http.NewRequest("PUT", mockURL, reader)
	if err != nil {
		require.NoError(c.T, err, "Failed to create request to mock server.")
	}
	mockRes, err := hc.Do(mockReq)
	if err != nil {
		require.NoError(c.T, err, "Failed to send request to mock server.")
	}
	// MockServer verification returns 202 on success and 406 on failure
	if mockRes.StatusCode != http.StatusAccepted {
		require.NoError(c.T, err,
			fmt.Sprintf("Mock server verification did not meet expectations and failed with status: %s", mockRes.Status))
	}
}
