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

type SleepSummaryDataField string
type SleepSummaryDataFields []SleepSummaryDataField

// String converts the slice of SleepSummaryDataFields into the string format expected by the API.
func (m SleepSummaryDataFields) String() string {
	v := make([]string, 0, len(m))
	for _, t := range m {
		v = append(v, string(t))
	}

	return strings.Join(v, ",")
}

const (
	SleepSummaryDataFieldNBREMEpisodes                  SleepSummaryDataField = "nb_rem_episodes"
	SleepSummaryDataFieldSleepEfficiency                SleepSummaryDataField = "sleep_efficiency"
	SleepSummaryDataFieldSleepLatency                   SleepSummaryDataField = "sleep_latency"
	SleepSummaryDataFieldTotalSleepTime                 SleepSummaryDataField = "total_sleep_time"
	SleepSummaryDataFieldTotalTimeInBed                 SleepSummaryDataField = "total_timeinbed"
	SleepSummaryDataFieldWakeupLatency                  SleepSummaryDataField = "wakeup_latency"
	SleepSummaryDataFieldWASO                           SleepSummaryDataField = "waso"
	SleepSummaryDataFieldApneaHyponeaIndex              SleepSummaryDataField = "apnea_hypopnea_index"
	SleepSummaryDataFieldBreathingDisturbancesIntensity SleepSummaryDataField = "breathing_disturbances_intensity"
	SleepSummaryDataFieldAsleepDuration                 SleepSummaryDataField = "asleepduration"
	SleepSummaryDataFieldDeepSleepDuration              SleepSummaryDataField = "deepsleepduration"
	SleepSummaryDataFieldDurationToSleep                SleepSummaryDataField = "durationtosleep"
	SleepSummaryDataFieldDurationToWakeup               SleepSummaryDataField = "durationtowakeup"
	SleepSummaryDataFieldHRAverage                      SleepSummaryDataField = "hr_average"
	SleepSummaryDataFieldHRMax                          SleepSummaryDataField = "hr_max"
	SleepSummaryDataFieldHRMin                          SleepSummaryDataField = "hr_min"
	SleepSummaryDataFieldLightSleepDuration             SleepSummaryDataField = "lightsleepduration"
	SleepSummaryDataFieldNightEvents                    SleepSummaryDataField = "night_events"
	SleepSummaryDataFieldOutOfBedCount                  SleepSummaryDataField = "out_of_bed_count"
	SleepSummaryDataFieldREMSleepDuration               SleepSummaryDataField = "remsleepduration"
	SleepSummaryDataFieldRRAverage                      SleepSummaryDataField = "rr_average"
	SleepSummaryDataFieldRRMax                          SleepSummaryDataField = "rr_max"
	SleepSummaryDataFieldRRMin                          SleepSummaryDataField = "rr_min"
	SleepSummaryDataFieldSleepScore                     SleepSummaryDataField = "sleep_score"
	SleepSummaryDataFieldSnoring                        SleepSummaryDataField = "snoring"
	SleepSummaryDataFieldSnoringEpisodeCount            SleepSummaryDataField = "snoringepisodecount"
	SleepSummaryDataFieldWakeUpCount                    SleepSummaryDataField = "wakeupcount"
	SleepSummaryDataFieldWakeUpDuration                 SleepSummaryDataField = "wakeupduration"
)

// SleepSummary is a summary fo sleep as defined by the Withings API.
type SleepSummary struct {
	Timezone  string `json:"timezone"`
	Model     int    `json:"model"`
	ModelID   int    `json:"model_id"`
	StartDate int    `json:"startdate"`
	EndDate   int    `json:"enddate"`
	Date      string `json:"date"`
	Created   int    `json:"created"`
	Modified  int    `json:"modified"`
	Data      struct {
		ApneaHypopneaIndex             *float64      `json:"apnea_hypopnea_index"`
		Asleepduration                 *float64      `json:"asleepduration"`
		BreathingDisturbancesIntensity *float64      `json:"breathing_disturbances_intensity"`
		DeepSleepDuration              *float64      `json:"deepsleepduration"`
		DurationtoSleep                *float64      `json:"durationtosleep"`
		DurationToWakeup               *float64      `json:"durationtowakeup"`
		HRAverage                      *float64      `json:"hr_average"`
		HRMax                          *float64      `json:"hr_max"`
		HRMin                          *float64      `json:"hr_min"`
		LightSleepDuration             *float64      `json:"lightsleepduration"`
		NBRemEpisodes                  *float64      `json:"nb_rem_episodes"`
		NightEvents                    []interface{} `json:"night_events"`
		OutOfBedCount                  *float64      `json:"out_of_bed_count"`
		REMSleepDuration               *float64      `json:"remsleepduration"`
		RrAverage                      *float64      `json:"rr_average"`
		RrMax                          *float64      `json:"rr_max"`
		RrMin                          *float64      `json:"rr_min"`
		SleepEfficiency                *float64      `json:"sleep_efficiency"`
		SleepLatency                   *float64      `json:"sleep_latency"`
		SleepScore                     *float64      `json:"sleep_score"`
		Snoring                        *float64      `json:"snoring"`
		SnoringEpisodeCount            *float64      `json:"snoringepisodecount"`
		TotalSleepTime                 *float64      `json:"total_sleep_time"`
		TotalTimeInBed                 *float64      `json:"total_timeinbed"`
		WakeupLatency                  *float64      `json:"wakeup_latency"`
		WakeupCount                    *float64      `json:"wakeupcount"`
		WakeupDuration                 *float64      `json:"wakeupduration"`
		WASO                           *float64      `json:"waso"`
	} `json:"data"`
}

// SleepSummaries is a slice of SleepSummary structs as defined by the Withings API.
type SleepSummaries []SleepSummary

// GetSleepSummaryResp is the response type returned by the Withings API for are request for sleep summary data.
type GetSleepSummaryResp struct {
	Status   int64               `json:"status"`
	APIError string              `json:"error"`
	Body     GetSleepSummaryBody `json:"body"`
}

// GetSleepSummaryBody is the body of the response returned by the Withings API for are request for sleep summary data.
type GetSleepSummaryBody struct {
	Series SleepSummaries `json:"series"`
	More   bool           `json:"more"`
	Offset int64          `json:"offset"`
}

// GetSleepSummaryParam contains the parameters needed to request sleep summary data.
type GetSleepSummaryParam struct {
	// Specifies the data fields that should be returned for each sleep summary.
	DataFields SleepSummaryDataFields

	// An offset value used for paging. The API response will return more with a 1 value if there are more pages
	// to retrieve. Along with this an offset value is provided. That value should be provided here on the next
	// request. See the Withings documentation for more information.
	Offset int64

	// Requests all data that was updated or created after this date. This is especially useful for data syncs
	// because it includes updated values which would not be included with StartDate and EndDate. If this value is
	// provided along with StartDate and EndDate, StartDate and EndDate will be ignored.
	LastUpdate time.Time
}

// UpdateQuery updates the query provided with the parameters of this param.
func (p *GetSleepSummaryParam) UpdateQuery(q url.Values) url.Values {
	// Constructing the query parameters based on the param provided.
	q.Set("action", APIActionGetSleepSummary)
	q.Set("data_fields", p.DataFields.String())
	if p.Offset > 0 {
		q.Set("offset", strconv.FormatInt(p.Offset, 10))
	}

	q.Set("lastupdate", strconv.FormatInt(p.LastUpdate.Unix(), 10))

	return q
}

// GetSleepSummary retrieves sleep summary data for the user represented by the token. Error will be non nil upon an internal
// or api error. If the API returned the error the response will contain the error.
// Due to an oddity of the Withings API it's possible you may receive an error due to failure to unmarshal the data.
// This happens when the scoped data fields are not found and the entry returns no data. Instead of an empty object
// the API returns an empty array. It's recommended to simply not provide any data fields.
func (c *Client) GetSleepSummary(ctx context.Context, token AccessToken, param GetSleepSummaryParam) (*GetSleepSummaryResp, error) {

	// Construct authorized request to request data from the API.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APISleepV2, nil)
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

	var mResp GetSleepSummaryResp
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
