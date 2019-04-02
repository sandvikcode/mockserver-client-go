package mock

// VerifyTimes checks how many times the Mock-Server was called for a given pattern
type VerifyTimes struct {
	AtLeast int `json:"atLeast,omitempty"`
	AtMost  int `json:"atMost,omitempty"`
}

// Verification checks requests have been received by MockServer
type Verification struct {
	Times *VerifyTimes `json:"times,omitempty"`
}
