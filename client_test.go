package withings_test

import (
	"context"
	"fmt"
	"github.com/jrmycanady/withings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

// client contains the test client configured on init.
var client *withings.Client

// demoToken contains the token generated for the demo user on init.
var demoToken *withings.AccessToken

// testingConfig contains the testing configuration loaded from the environmental variables.
var testingConfig = struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}{}

func init() {
	loadTestConfigFromEnv()
}

func loadTestConfigFromEnv() {
	testingConfig.ClientID = os.Getenv("GO_WITHINGS_TEST_CLIENT_ID")
	testingConfig.ClientSecret = os.Getenv("GO_WITHINGS_TEST_CLIENT_SECRET")
	testingConfig.RedirectURL = os.Getenv("GO_WITHINGS_TEST_REDIRECT_URL")
}

func TestClient_SetHTTPClientTimeout_Option(t *testing.T) {
	c := withings.NewClient(testingConfig.ClientID, testingConfig.ClientSecret, url.URL{}, withings.SetHTTPClientTimeout(10*time.Second))
	assert.Equal(t, 10*time.Second, c.HttpClient.Timeout)
}

func TestClient_WithDemoMode_Option(t *testing.T) {
	redirectURL, err := url.Parse(testingConfig.RedirectURL)
	require.Nil(t, err)
	require.NotNil(t, redirectURL)
	c := withings.NewClient(testingConfig.ClientID, testingConfig.ClientSecret, *redirectURL, withings.WithDemoMode())

	authURL, _, err := c.GetUserAuthRequestURL([]string{}, "")
	require.Nil(t, err)
	require.NotNil(t, authURL)
	assert.Contains(t, authURL.String(), "demo")

}

func TestClient_WithSkipSSLVerify_Option(t *testing.T) {
	c := withings.NewClient(testingConfig.ClientID, testingConfig.ClientSecret, url.URL{}, withings.WithSkipSSLVerify())

	switch transport := c.HttpClient.Transport.(type) {
	case *http.Transport:
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	default:
		require.Fail(t, "http client transport is not an *http.Transport")
	}
}

func TestClient_GetAuthenticationRequestURL(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		scopes      []string
		state       string
		redirectURL string
	}{
		"All scopes and no state": {
			scopes:      []string{withings.ScopeUserActivity, withings.ScopeUserMetrics},
			state:       "",
			redirectURL: testingConfig.RedirectURL,
		},
		"No scopes and no state": {
			scopes:      []string{},
			state:       "",
			redirectURL: testingConfig.RedirectURL,
		},
		"All scopes and provided state": {
			scopes:      []string{withings.ScopeUserActivity, withings.ScopeUserMetrics},
			state:       "UNIQUESTATECHECK",
			redirectURL: testingConfig.RedirectURL,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Building and validating the URL before testing the actual method.
			redirectURL, err := url.Parse(test.redirectURL)
			require.Nil(t, err)
			require.NotNil(t, redirectURL)

			c := withings.NewClient(testingConfig.ClientID, testingConfig.ClientSecret, *redirectURL, withings.SetHTTPClientTimeout(10*time.Second))

			authURL, state, err := c.GetUserAuthRequestURL(test.scopes, test.state)
			require.Nil(t, err)

			assert.NotEmpty(t, state)

			for _, s := range test.scopes {
				assert.Contains(t, authURL.String(), s)
			}

			assert.Contains(t, authURL.String(), url.QueryEscape(test.redirectURL))

			if test.state != "" {
				assert.Contains(t, authURL.String(), test.state)
				assert.Equal(t, test.state, state)
			}
		})
	}

}

// getDemoClient builds a new client and retrieves an access token for the demo user.
func getDemoClient() (*withings.Client, *withings.AccessToken, error) {
	// Building and validating the URL before testing the actual method.
	redirectURL, err := url.Parse(testingConfig.RedirectURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to  build url: %s", err)
	}

	c := withings.NewClient(testingConfig.ClientID, testingConfig.ClientSecret, *redirectURL, withings.SetHTTPClientTimeout(10*time.Second), withings.WithDemoMode())

	accessToken, err := c.GetDemoAccessToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get access token: %s", err)
	}
	if accessToken == nil {
		return nil, nil, fmt.Errorf("access token was nil: %s", err)
	}
	if accessToken.Status != 0 {
		return nil, nil, fmt.Errorf("received error %d when retrieving demo access token", accessToken.Status)
	}

	return c, &accessToken.AccessToken, nil
}

func init() {
	// Generating the demo client and token. All tests should check these before proceeding with testing.
	var err error
	client, demoToken, err = getDemoClient()
	if err != nil {
		panic(err)
	}

}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestClient_GetMeasure(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param                 withings.GetMeasureParam
		status                int64
		expectedFirstResult   withings.MeasureGroup
		expectedGroupCount    int
		expectedTotalMeasures int
	}{
		"Retrieve unbound weights only": {
			param: withings.GetMeasureParam{
				MeasurementTypes: withings.MeasureTypes{withings.MeasureTypeWeightKilogram},
			},
		},
		//"Retrieve unbound all measurements": {
		//	param: withings.GetMeasureParam{},
		//},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetMeasure(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}

func TestClient_GetActivity(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param                 withings.GetActivityParam
		status                int64
		expectedFirstResult   withings.Activity
		expectedGroupCount    int
		expectedTotalMeasures int
	}{
		"Retrieve unbound": {
			param: withings.GetActivityParam{
				LastUpdate: time.Now().Add(-24 * time.Hour),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetActivity(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}

func TestClient_GetIntraDayActivity(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.GetIntraDayActivityParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.GetIntraDayActivityParam{
				DataFields: withings.IntraDayActivityFields{
					withings.IntraDayActivityFieldCalories,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetIntraDayActivity(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
			assert.Greater(t, len(resp.Body.Series), 0)
		})
	}
}

func TestClient_GetWorkout(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.GetWorkoutParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.GetWorkoutParam{
				DataFields: withings.WorkoutDataFields{
					withings.WorkoutDataFieldCalories,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetWorkout(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
			assert.Greater(t, len(resp.Body.Series), 0)
		})
	}

}

func TestClient_GetHeartData(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.GetHeartListParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.GetHeartListParam{
				StartDate: timePtr(time.Unix(1594159644-1000, 0)),
				EndDate:   timePtr(time.Now()),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetHeartList(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}

func TestClient_GetSleep(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.GetSleepParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.GetSleepParam{
				StartDate: time.Now().Add(-300 * time.Hour),
				EndDate:   time.Now(),
				DataFields: withings.SleepDataFields{
					withings.SleepDataFieldRR,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetSleep(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
			assert.Greater(t, len(resp.Body.Series), 0)
		})
	}

}

func TestClient_GetSleepSummary(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.GetSleepSummaryParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.GetSleepSummaryParam{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.GetSleepSummary(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
			assert.Greater(t, len(resp.Body.Series), 0)
		})
	}

}

func TestClient_ListNotification(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param  withings.ListNotificationParam
		status int64
	}{
		"Retrieve unbound": {
			param: withings.ListNotificationParam{
				Appli: 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.ListNotification(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}

func TestClient_GetMeasureWithRefreshedToken(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param                 withings.GetMeasureParam
		status                int64
		expectedFirstResult   withings.MeasureGroup
		expectedGroupCount    int
		expectedTotalMeasures int
	}{
		"Retrieve unbound weights only": {
			param: withings.GetMeasureParam{
				MeasurementTypes: withings.MeasureTypes{withings.MeasureTypeWeightKilogram},
			},
		},
		//"Retrieve unbound all measurements": {
		//	param: withings.GetMeasureParam{},
		//},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			refreshResult, err := client.RefreshAccessToken(*demoToken)
			require.Nil(t, err)
			demoToken = &refreshResult.AccessToken
			resp, err := client.GetMeasure(context.Background(), *demoToken, test.param)
			require.Nil(t, err)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}

func TestAuthorizedUser_GetMeasure(t *testing.T) {
	t.Parallel()

	// Verify init succeeded.
	require.NotNil(t, client)
	require.NotNil(t, demoToken)

	tests := map[string]struct {
		param                 withings.GetMeasureParam
		status                int64
		expectedFirstResult   withings.MeasureGroup
		expectedGroupCount    int
		expectedTotalMeasures int
	}{
		"Retrieve unbound weights only": {
			param: withings.GetMeasureParam{
				MeasurementTypes: withings.MeasureTypes{withings.MeasureTypeWeightKilogram},
			},
		},
		//"Retrieve unbound all measurements": {
		//	param: withings.GetMeasureParam{},
		//},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := client.NewAuthorizedUser(*demoToken)

			resp, token, err := u.GetMeasure(context.Background(), test.param)
			require.Nil(t, err)
			require.Nil(t, token)
			require.Equal(t, int64(0), resp.Status)
		})
	}

}
