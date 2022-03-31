package withings

import (
	"context"
	"sync"
	"time"
)

// AuthorizedUser is a user that has granted the client access to their data via an access token.
type AuthorizedUser struct {
	c *Client
	t *AccessToken
	sync.Mutex
}

func (c *Client) NewAuthorizedUser(t AccessToken) *AuthorizedUser {
	return &AuthorizedUser{
		c: c,
		t: &t,
	}
}

// checkToken checks if the token is still valid and requests a new token if needed. If a new token
// is obtained it is returned.
func (a *AuthorizedUser) checkToken() (*AccessTokenResponse, error) {

	// Locking for the entire life of the call to prevent any other attempts with the token.
	a.Lock()
	defer a.Unlock()

	expAt := a.t.ExpiresAt.Add(-10 * time.Second)

	if time.Now().After(expAt) {
		tokenResp, err := a.c.RefreshAccessToken(*a.t)
		if err != nil {
			return tokenResp, err
		}
		a.t = &tokenResp.AccessToken
		return tokenResp, nil
	}

	return nil, nil
}

// GetMeasure returns the measures for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetMeasure(ctx context.Context, param GetMeasureParam) (*GetMeasureResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetMeasure(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetIntraDayActivity returns the intra day activities for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetIntraDayActivity(ctx context.Context, param GetIntraDayActivityParam) (*GetIntraDayActivityResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetIntraDayActivity(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetActivity returns the activities for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetActivity(ctx context.Context, param GetActivityParam) (*GetActivityResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetActivity(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetHeartList returns the Heart Data for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetHeartList(ctx context.Context, param GetHeartListParam) (*GetHeartResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetHeartList(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetHeartHighFrequencyData returns the Heart Data for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetHeartHighFrequencyData(ctx context.Context, param GetHeartHighFrequencyDataParam) (*GetHeartHighFrequencyDataResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetHeartHighFrequencyData(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetSleep returns the Sleep data for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetSleep(ctx context.Context, param GetSleepParam) (*GetSleepResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetSleep(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetSleepSummary returns the SleepSummary data for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetSleepSummary(ctx context.Context, param GetSleepSummaryParam) (*GetSleepSummaryResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetSleepSummary(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}

// GetWorkout returns the Workout data for the AuthorizedUser based on the param provided. If a new token had to be created
// it will be non nil.
func (a *AuthorizedUser) GetWorkout(ctx context.Context, param GetWorkoutParam) (*GetWorkoutResp, *AccessToken, error) {
	tokenResp, err := a.checkToken()
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.c.GetWorkout(ctx, *a.t, param)

	if tokenResp != nil {
		return resp, &tokenResp.AccessToken, err
	}

	return resp, nil, err
}
