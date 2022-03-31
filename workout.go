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

type WorkoutDataField string
type WorkoutDataFields []WorkoutDataField

// String converts the slice of ActivityDataFields into the string format expected by the API.
func (m WorkoutDataFields) String() string {
	v := make([]string, 0, len(m))
	for _, t := range m {
		v = append(v, string(t))
	}

	return strings.Join(v, ",")
}

const (
	WorkoutDataFieldCalories          WorkoutDataField = "calories"
	WorkoutDataFieldIntensity         WorkoutDataField = "intensity"
	WorkoutDataFieldManualDistance    WorkoutDataField = "manual_distance"
	WorkoutDataFieldManualCalories    WorkoutDataField = "manual_calories"
	WorkoutDataFieldHRAverage         WorkoutDataField = "hr_average"
	WorkoutDataFieldHRMin             WorkoutDataField = "hr_min"
	WorkoutDataFieldHRMAx             WorkoutDataField = "hr_max"
	WorkoutDataFieldHRZone0           WorkoutDataField = "hr_zone_0"
	WorkoutDataFieldHRZone1           WorkoutDataField = "hr_zone_1"
	WorkoutDataFieldHRZone2           WorkoutDataField = "hr_zone_2"
	WorkoutDataFieldHRZone3           WorkoutDataField = "hr_zone_3"
	WorkoutDataFieldPauseDuration     WorkoutDataField = "pause_duration"
	WorkoutDataFieldAlgoPauseDuration WorkoutDataField = "algo_pause_duration"
	WorkoutDataFieldSPO2Average       WorkoutDataField = "spo2_average"
	WorkoutDataFieldSteps             WorkoutDataField = "steps"
	WorkoutDataFieldDistance          WorkoutDataField = "distance"
	WorkoutDataFieldElevation         WorkoutDataField = "elevation"
	WorkoutDataFieldPoolLaps          WorkoutDataField = "pool_laps"
	WorkoutDataFieldStrokes           WorkoutDataField = "strokes"
	WorkoutDataFieldPoolLength        WorkoutDataField = "pool_length"
)

// Workout is a workout as defined by the Withings API.
type Workout struct {
	Category  int         `json:"category"`
	Timezone  string      `json:"timezone"`
	Model     int         `json:"model"`
	Attrib    int         `json:"attrib"`
	StartDate int         `json:"startdate"`
	EndDate   int         `json:"enddate"`
	Date      string      `json:"date"`
	Modified  int         `json:"modified"`
	DeviceID  string      `json:"deviceid"`
	Data      WorkoutData `json:"data"`
}

// Workouts is a slice of Workout structs.
type Workouts []Workout

// WorkoutData is the data of a workout as defined by the Withings API.
type WorkoutData struct {
	AlgoPauseDuration *float64 `json:"algo_pause_duration"`
	Calories          *float64 `json:"calories"`
	Distance          *float64 `json:"distance"`
	Elevation         *float64 `json:"elevation"`
	HrAverage         *float64 `json:"hr_average"`
	HrMax             *float64 `json:"hr_max"`
	HrMin             *float64 `json:"hr_min"`
	HrZone0           *float64 `json:"hr_zone_0"`
	HrZone1           *float64 `json:"hr_zone_1"`
	HrZone2           *float64 `json:"hr_zone_2"`
	HrZone3           *float64 `json:"hr_zone_3"`
	Intensity         *float64 `json:"intensity"`
	ManualCalories    *float64 `json:"manual_calories"`
	ManualDistance    *float64 `json:"manual_distance"`
	PauseDuration     *float64 `json:"pause_duration"`
	PoolLaps          *float64 `json:"pool_laps"`
	PoolLength        *float64 `json:"pool_length"`
	Spo2Average       *float64 `json:"spo2_average"`
	Steps             *float64 `json:"steps"`
	Strokes           *float64 `json:"strokes"`
}

// GetWorkoutResp is the response type returned by the Withings API for are request for workout data.
type GetWorkoutResp struct {
	Status   int64          `json:"status"`
	APIError string         `json:"error"`
	Body     GetWorkoutBody `json:"body"`
}

// GetWorkoutBody is the body of the response returned by the Withings API for are request for workout data.
type GetWorkoutBody struct {
	Series Workouts `json:"series"`
	More   bool     `json:"more"`
	Offset int64    `json:"offset"`
}

// GetWorkoutParam contains the parameters needed to request workouts.
type GetWorkoutParam struct {
	// Specifies the data fields that should be returned for each workout.
	DataFields WorkoutDataFields

	// The start of the window of measurements to retrieve. This value is ignored if LastUpdate is provided.
	StartDate *time.Time

	// The end of the window of the measurements to retrieve. This value is ignored if LastUpdate is provided.
	EndDate *time.Time

	// An offset value used for paging. The API response will return more with a 1 value if there are more pages
	// to retrieve. Along with this an offset value is provided. That value should be provided here on the next
	// request. See the Withings documentation for more information.
	Offset int64

	// Requests all data that was updated or created after this date. This is especially useful for data syncs
	// because it includes updated values which would not be included with StartDate and EndDate. If this value is
	// provided along with StartDate and EndDate, StartDate and EndDate will be ignored.
	LastUpdate *time.Time
}

// UpdateQuery updates the query provided with the parameters of this param.
func (p *GetWorkoutParam) UpdateQuery(q url.Values) url.Values {
	// Constructing the query parameters based on the param provided.
	q.Set("action", APIActionGetWorkout)
	q.Set("data_fields", p.DataFields.String())
	if p.Offset > 0 {
		q.Set("offset", strconv.FormatInt(p.Offset, 10))
	}
	switch p.LastUpdate {
	case nil:
		if p.StartDate != nil {
			q.Set("startdate", strconv.FormatInt(p.StartDate.Unix(), 10))
		}
		if p.EndDate != nil {
			q.Set("enddate", strconv.FormatInt(p.EndDate.Unix(), 10))
		}
	default:
		q.Set("lastupdate", strconv.FormatInt(p.LastUpdate.Unix(), 10))
	}

	return q
}

// GetWorkout retrieves workouts for the user represented by the token. Error will be non nil upon an internal
// or api error. If the API returned the error the response will contain the error.
func (c *Client) GetWorkout(ctx context.Context, token AccessToken, param GetWorkoutParam) (*GetWorkoutResp, error) {

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

	var mResp GetWorkoutResp
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
