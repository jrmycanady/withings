# Withings


```go

// Create a new client with your parameters as supplied by Withings when registring your app.
client : = withings.NewClient(clientID, clientSecret, redirectURL, with...)


// Obtaining access to a users data.
client.GetUserAuthRequestURL()

// User access the URL and is redirected to you redirectURL... 

c.GetUserAccesToken(authCode)

// Now we have an access token. We can then use it to retrieve data by creating a new
// userAccess. The  userAccess allows us to perform requests and will automatically update
// any tokens as needed. You can

c := withings.NewClient()

c.UserAuthURL()

u := c.AuthUser(code)

u := c.UserFromAccessToken(AccessToken, WithAutoRefresh())

// Data contains the data from the request.
// AccessToken is provided if the access token changed to prevent expiry.
// Error is provided if any error happend. It will be of type APIError if the error was api related.
data, accessToken, err := u.GetMeasures(param)

u.GetActivity
u.GetIntraDayActivity
u.GetWorkouts
u.HeartList
u.HeartGet
u.SleepGet
u.SleepGetSummary
u.NotifyGet
u.NotifyList
u.NotifyRevoke
u.NotifySubsribe
u.NotifyUpdate
u.GetDevices
u.GetGoals


```

### Test Env 

|Name|Description|
|----|-----------|
|GO_WITHINGS_TEST_CLIENT_ID|The ClientID to use when testing.|
|GO_WITHINGS_TEST_CLIENT_SECRET|The ClientSecret to use when testing.|
|GO_WITHINGS_TEST_REDIRECT_URL|The RedirectURL to use when testing.|