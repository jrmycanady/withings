package withings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Device is device as defined by the Withings API.
type Device struct {
	Type            string `json:"type"`
	Model           string `json:"model"`
	ModelID         int    `json:"model_id"`
	Battery         string `json:"battery"`
	DeviceID        string `json:"deviceid"`
	HashDeviceID    string `json:"hash_deviceid"`
	Timezone        string `json:"timezone"`
	LastSessionDate int    `json:"last_session_date"`
}

// Devices is a slice of Device structs.
type Devices []Device

// GetUserDeviceResp is the response type returned by the Withings API for are request for user devices.
type GetUserDeviceResp struct {
	Status   int64             `json:"status"`
	APIError string            `json:"error"`
	Body     GetUserDeviceBody `json:"body"`
}

// GetUserDeviceBody is the body of the response returned by the Withings API for are request for user devices.
type GetUserDeviceBody struct {
	Devices Devices `json:"devices"`
	More    bool    `json:"more"`
	Offset  int64   `json:"offset"`
}

// GetUserDevice retrieves devices for the user represented by the token. Error will be non nil upon an internal
// or api error. If the API returned the error the response will contain the error.
func (c *Client) GetUserDevice(ctx context.Context, token AccessToken) (*GetUserDeviceResp, error) {

	// Construct authorized request to request data from the API.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIUser, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	q := req.URL.Query()
	q.Set("action", APIActionUserGetDevice)
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body of request: %w", err)
	}

	var mResp GetUserDeviceResp
	if err = json.Unmarshal(body, &mResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	switch mResp.Status {
	case 0:
		return &mResp, nil
	default:
		return &mResp, fmt.Errorf("api returned an error: %s", mResp.APIError)
	}
}
