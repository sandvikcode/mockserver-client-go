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

type Client struct {
	// testing type T
	T *testing.T
	// mock server base address
	BaseURL string
}

func (hms *Client) AddExpectation(exp *Expectation) {
	msg, err := json.Marshal(exp)
	if err != nil {
		require.NoError(hms.T, err,
			"Failed to serialize mock server expectation.")
	}

	hms.callMock("expectation", string(msg))
}

// Clear some expectations for a given path in Mock-Server
// TODO: refactor with a clear model
// TODO: this could be part of expectations.go? http://www.mock-server.com/mock_server/clearing_and_resetting.html
func (hms *Client) Clear(path string) {
	mockReqBody := fmt.Sprintf(`
			{
				"path": "%s"
			}
			`, path)
	hms.callMock("clear", mockReqBody)
}

// VerifyMinMax Mock-Server was called at least and at most N number of times. Note: -1 is infinite number of times.
// TODO: refactor with a clear model
// TODO: Move to verifications.go For the model we will have VerifyCount and VerifySequence, see http://www.mock-server.com/mock_server/verification.html
func (hms *Client) VerifyMinMax(path string, atLeast int, atMost int) {
	mockReqBody := fmt.Sprintf(`
			{
				"httpRequest": {
					"path": "%s"
				},
				"times": {
					"atLeast": %d,
					"atMost": %d
				}
			}
			`, path, atLeast, atMost)

	hms.callMock("mockserver/verify", mockReqBody)
}

// Reset the entire Mock-Server, clearing all state
// TODO: This should be part of client as it clears everything..verifications, expectations etc
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
	// Mock-server verification returns 202 on success and 406 on failure
	if mockRes.StatusCode != http.StatusAccepted {
		require.NoError(hms.T, err,
			fmt.Sprintf("Mock server verification did not meet expectations and failed with status: %s", mockRes.Status))
	}
}