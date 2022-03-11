package withings

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var codeErrorRegex = regexp.MustCompile(`.*code=(.*)&.*`)
var csrfRegex = regexp.MustCompile(`.*csrf_token.*value="(.*)"`)

// GetDemoAccessToken creates an auth request for the demo user and obtains a valid access token. This is used
// for testing purposes.
func (c *Client) GetDemoAccessToken() (*AccessTokenResponse, error) {
	authReqURL, _, err := c.GetUserAuthRequestURL([]string{ScopeUserMetrics, ScopeUserActivity}, "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate request URL: %s", err)
	}

	// Perform request to build the form.
	getReq, err := http.NewRequest(http.MethodGet, authReqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build get request: %s", err)
	}
	getResp, err := c.HttpClient.Do(getReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get request: %s", err)
	}
	defer getResp.Body.Close()
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read get body: %s", err)
	}

	if getResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status of %d", getResp.StatusCode)
	}

	results := csrfRegex.FindSubmatch(getBody)
	if len(results) != 2 {
		return nil, fmt.Errorf("failed to find csrf in page")
	}
	csrf := string(results[1])

	formData := url.Values{}
	formData.Set("authorized", "1")
	formData.Set("csrf_token", csrf)

	postReq, err := http.NewRequest(http.MethodPost, authReqURL.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to build post request: %s", err)
	}
	postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	postResp, err := c.HttpClient.Do(postReq)
	var code string
	if err != nil {
		// We expect a possible error here if the redirect URI does not actually exist. Check for that case.
		if !strings.Contains(err.Error(), "no such host") {
			return nil, fmt.Errorf("failed to execute post request: %s", err)
		}

		results := codeErrorRegex.FindStringSubmatch(err.Error())
		if len(results) != 2 {
			return nil, fmt.Errorf("failed to find code in response error [%s]", err.Error())
		}
		code = results[1]

	} else {
		defer postResp.Body.Close()

		// We expect it to fail and just want to look at the redirected URL for the code.
		code = postResp.Request.URL.Query().Get("code")
	}

	fmt.Println(code)

	acessToken, err := c.GetUserAccessToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token with code: %s", err)
	}

	fmt.Println(acessToken)

	// The request was accepted so issue the request with the proper form data.
	return acessToken, nil
}
