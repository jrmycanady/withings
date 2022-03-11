package withings

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ScopeUserActivity = "user.activity"
	ScopeUserMetrics  = "user.metrics"
)

const (
	APIActionGetMeasure  = "getmeas"
	APIActionGetActivity = "getactivity"
)

type Client struct {

	// Contains the Client ID that is associated with this client. This Client ID is provided when a client application
	// is registered with withings. See https://developer.withings.com
	clientID string

	// Contains the secret that is associated with the Client ID.
	clientSecret string

	// Contains the HTTP client that will be used for client level HTTP calls.
	HttpClient *http.Client

	// Denotes if certificates verification should be skipped.
	skipCertificateVerification bool

	// Denotes the default timeout duration for http clients.
	httpClientTimeout time.Duration

	// Contains the URL that the Withings API redirects to during authentication actions.
	redirectURL url.URL

	// Denotes the client should run in demo mode.
	demoMode bool
}

func NewClient(clientID string, clientSecret string, redirectURL url.URL, opts ...ClientOption) *Client {
	c := &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}

	// Apply options.
	for _, opt := range opts {
		opt(c)
	}

	// Building the default http client with specified values.
	c.HttpClient = &http.Client{
		Timeout: c.httpClientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.skipCertificateVerification,
			},
		},
	}

	return c
}

// ClientOption is an option that can be applied to a client.
type ClientOption func(client *Client)

// WithSkipSSLVerify configures the client to skip the verification of all SSL certificates.
func WithSkipSSLVerify() ClientOption {
	return func(c *Client) {
		c.skipCertificateVerification = true
	}
}

// SetHTTPClientTimeout configures the default timeout value used for all HTTP clients.
func SetHTTPClientTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClientTimeout = timeout
	}
}

// WithDemoMode configures the client to run in API demo mode. This is only needed for testing.
func WithDemoMode() ClientOption {
	return func(c *Client) {
		c.demoMode = true
	}
}

// GetUserAuthRequestURL generates the URL that a user must access to grant this client's access to their Withings
// data. The scope of access is determined by the scopes provided. After successful granting, the Withings API will
// redirect the user to the redirectURL specified. This URL must be set to the same URL base as the value set for the
// Callback URI configured when registering the client application.
//
// The API also accepts a state value that is provided back to validate the redirect wasn't spoofed. The state can be
// provided but if empty a randomly generated Base64 string will be generated.
func (c *Client) GetUserAuthRequestURL(scopes []string, state string) (authRequestURL *url.URL, expectedState string, err error) {

	// Building base request.
	authRequestURL, err = url.Parse(APIPathUserAuthorize)
	if err != nil {
		// This must never fail. Panic here so tests fail hard and fast
		panic(err)
	}
	query := authRequestURL.Query()

	// Generating state if needed.
	if state == "" {
		v := make([]byte, 32)
		_, err = io.ReadFull(rand.Reader, v[:])
		if err != nil {
			return authRequestURL, "", fmt.Errorf("failed to generate state value: %w", err)
		}
		state = base64.URLEncoding.EncodeToString(v)
	}

	// Set per the API spec.
	query.Set("response_type", "code")

	query.Set("client_id", c.clientID)
	query.Set("state", state)
	query.Set("scope", strings.Join(scopes, ","))
	query.Set("redirect_uri", c.redirectURL.String())

	// Configuring for the demo mode if needed.
	if c.demoMode {
		query.Set("mode", "demo")
	}

	authRequestURL.RawQuery = query.Encode()

	return authRequestURL, state, nil
}

type AccessTokenResponse struct {
	Status      int64       `json:"status"`
	AccessToken AccessToken `json:"body"`
}
type AccessToken struct {
	UserID       int64  `json:"userid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	CSRFToken    string `json:"csrf_token"`
	TokenType    string `json:"token_type"`
}

// GetUserAccessToken retrieves a new user access token using the AuthCode provided. The authCode is provided by the
// user visiting the URL provided by GetUserAuthenticationRequestURL and allowing access. The redirectURL provided
// must match the URL provided during generation of the authCode.
func (c *Client) GetUserAccessToken(authCode string) (*AccessTokenResponse, error) {

	// Building required form data for the request.
	formData := url.Values{}
	formData.Set("action", "requesttoken")
	formData.Set("client_id", c.clientID)
	formData.Set("client_secret", c.clientSecret)
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCode)
	formData.Set("redirect_uri", c.redirectURL.String())

	req, err := http.NewRequest(http.MethodPost, APIPathUserAccessToken, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %s", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var accessToken AccessTokenResponse
		if err = json.Unmarshal(body, &accessToken); err != nil {
			return nil, fmt.Errorf("failed to parse response: %s", err)
		}
		return &accessToken, nil
	default:
		return nil, fmt.Errorf("failed with API error")
	}
}
