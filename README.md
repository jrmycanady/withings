# Withings Client

This module is a go client for the public Withings API. It provides reasonable coverage of the API to allow accessing
user data. 

## Supported 

* Measurements (Weight, Blood Pressure, etc)
* Activities
* Intra Day Activities
* Workouts
* Heart Rate Data
* High Frequency Heart Rate Data
* Sleep
* Sleep Summary

## Installation
> go get github.com/jrmycanady/withings@latest

## Basic Usage

The withings API utilizes OAuth 2.0 to provide access to the API. If you are not familiar with OAuth 2.0 and Access/Refresh
tokens you should review the [Withings API documentation](https://developer.withings.com/developer-guide/v3/integration-guide/public-health-data-api/get-access/oauth-web-flow). Once you have registered 
your application and have the codes/tokens you can then proceed.

1. Create a new client configured with your client ID, secret, and redirect URL. Note, you may provide various client options to modify the client.
```go
c := withings.NewClient("id", "secret", "redirect url")
```

2. Users that have not granted your application access will need to so. This is done by generating an authorization URL the user navigates to via a browser. Once they click allow, they will be redirected back to your redirection URL with the code you need to finish generating the token that will give your application access. You must specify a scope that denotes what data you want to access. In this example we are specifying all measurements and activities. Scopes can be found [here](https://developer.withings.com/developer-guide/v3/data-api/all-available-health-data).

```go
authURL, state, err := c.GetUserAuthRequestURL([]string{withings.ScopeUserActivity, withings.ScopeUserMetrics}, "")
```

3. Once you have the code you can generate an access and refresh token for the user.

```go
token, err := c.GetUserAccessToken('code')
```

4. You may now use the token to access the users data. The module provides two ways to do this. The first is a direct call with you providing the specific token to use. If the token is expires it will result in an API error. The second is to generate a AuthorizedUser from the token. You may then use the user to perform the same data requests. The difference is the user will automatically update the access token if it's about to expire.
```go
// Direct method.
resp, err := c.GetWorkout(context.Background(), token, param)

// User Method
u := withings.NewAuthorizedUser(token)
resp, err := u.GetWorkout(context.Background(), param)
```

## Measures Access Methods

By default, the module returns the data in the format provided by the Withings API allowing you to work with it any way you like. For convince some types also have access methods to aid in accessing data. The MeasureGroups type allows for retrieving all measurements of one type from the dataset returned. 

```go
param := withings.GetMeasureParam{
    MeasurementTypes: withings.MeasureTypes{withings.MeasureTypeWeightKilogram},
}
resp, err := u.GetMeasure(context.Background(), param)
if err != nil {
	panic(err)
}

// Obtaining all the weight measurements across all measurement groups.
weights := resp.Weights()
```

### Test Env 

|Name|Description|
|----|-----------|
|GO_WITHINGS_TEST_CLIENT_ID|The ClientID to use when testing.|
|GO_WITHINGS_TEST_CLIENT_SECRET|The ClientSecret to use when testing.|
|GO_WITHINGS_TEST_REDIRECT_URL|The RedirectURL to use when testing.|