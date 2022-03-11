package withings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ActivityDataField string
type ActivityDataFields []ActivityDataField

// String converts the slice of ActivityDataFields into the string format expected by the API.
func (m ActivityDataFields) String() string {
	v := make([]string, 0, len(m))
	for _, t := range m {
		v = append(v, string(t))
	}

	return strings.Join(v, ",")
}

const (
	ActivityDataFieldSteps         ActivityDataField = "steps"
	ActivityDataFieldDistance      ActivityDataField = "distance"
	ActivityDataFieldElevation     ActivityDataField = "elevation"
	ActivityDataFieldSoft          ActivityDataField = "soft"
	ActivityDataFieldModerate      ActivityDataField = "moderate"
	ActivityDataFieldIntense       ActivityDataField = "intense"
	ActivityDataFieldActive        ActivityDataField = "active"
	ActivityDataFieldCalories      ActivityDataField = "calories"
	ActivityDataFieldTotalCalories ActivityDataField = "totalcalories"
	ActivityDataFieldHRAverage     ActivityDataField = "hr_average"
	ActivityDataFieldHRMin         ActivityDataField = "hr_min"
	ActivityDataFieldHRMAx         ActivityDataField = "hr_max"
	ActivityDataFieldHRZone0       ActivityDataField = "hr_zone_0"
	ActivityDataFieldHRZone1       ActivityDataField = "hr_zone_1"
	ActivityDataFieldHRZone2       ActivityDataField = "hr_zone_2"
	ActivityDataFieldHRZone3       ActivityDataField = "hr_zone_3"
)

// Activity is an activity as defined by the Withings API.
type Activity struct {
	Date          string  `json:"date"`
	Timezone      string  `json:"timezone"`
	DeviceID      string  `json:"deviceid"`
	HashDeviceID  string  `json:"hash_deviceid"`
	Brand         float64 `json:"brand"`
	IsTracker     bool    `json:"is_tracker"`
	Steps         float64 `json:"steps"`
	Distance      float64 `json:"distance"`
	Elevation     float64 `json:"elevation"`
	Soft          float64 `json:"soft"`
	Moderate      float64 `json:"moderate"`
	Intense       float64 `json:"intense"`
	Active        float64 `json:"active"`
	Calories      float64 `json:"calories"`
	TotalCalories float64 `json:"totalcalories"`
	HrAverage     float64 `json:"hr_average"`
	HrMin         float64 `json:"hr_min"`
	HrMax         float64 `json:"hr_max"`
	HrZone0       float64 `json:"hr_zone_0"`
	HrZone1       float64 `json:"hr_zone_1"`
	HrZone2       float64 `json:"hr_zone_2"`
	HrZone3       float64 `json:"hr_zone_3"`
}

// Activities is a slice of Activity structs as defined by the Withings API.
type Activities []Activity

// GetActivityResp is the response type returned by the Withings API for are request for activity data.
type GetActivityResp struct {
	Status   int64           `json:"status"`
	APIError string          `json:"error"`
	Body     GetActivityBody `json:"body"`
}

// GetActivityBody is the body of the response returned by the Withings API for are request for activity data.
type GetActivityBody struct {
	Activities Activities `json:"activities"`
	More       bool       `json:"more"`
	Offset     int64      `json:"offset"`
}

// GetActivityParam contains the parameters needed to request activities.
type GetActivityParam struct {
	// An offset value used for paging. The API response will return more with a 1 value if there are more pages
	// to retrieve. Along with this an offset value is provided. That value should be provided here on the next
	// request. See the Withings documentation for more information.
	Offset int64

	// Specifies the data fields that should be returned for each activity.
	DataFields ActivityDataFields

	// Requests all data that was updated or created after this date. This is especially useful for data syncs
	// because it includes updated values which would not be included with StartDate and EndDate. If this value is
	// provided along with StartDate and EndDate, StartDate and EndDate will be ignored.
	LastUpdate time.Time
}

// UpdateQuery updates the query provided with the parameters of this param.
func (p *GetActivityParam) UpdateQuery(q url.Values) url.Values {
	// Constructing the query parameters based on the param provided.
	q.Set("action", APIActionGetActivity)
	if len(p.DataFields) > 0 {
		q.Set("data_fields", p.DataFields.String())
	}
	if p.Offset > 0 {
		q.Set("offset", strconv.FormatInt(p.Offset, 10))
	}

	q.Set("lastupdate", strconv.FormatInt(p.LastUpdate.Unix(), 10))

	return q
}

// GetActivity retrieves activities for the user represented by the token. Error will be non nil upon an internal
// or api error. If the API returned the error the response will contain the error.
func (c *Client) GetActivity(ctx context.Context, token AccessToken, param GetActivityParam) (*GetActivityResp, error) {

	// Construct authorized request to request data from the API.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIPathGetV2Measure, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	// Updating the query with the parameters generated by the param provided.
	req.URL.RawQuery = param.UpdateQuery(req.URL.Query()).Encode()

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

	var mResp GetActivityResp
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
