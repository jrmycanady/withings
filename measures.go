package withings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MeasureCategory int64

const (
	MeasureCategoryReal      MeasureCategory = 1
	MeasureCategoryObjective MeasureCategory = 2
)

type MeasureType int64
type MeasureTypes []MeasureType

// String converts the slice of MeasureTypes into the string format expected by the API.
func (m MeasureTypes) String() string {
	v := make([]string, 0, len(m))
	for _, t := range m {
		v = append(v, strconv.FormatInt(int64(t), 10))
	}

	return strings.Join(v, ",")
}

const (
	MeasureTypeWeightKilogram                  MeasureType = 1
	MeasureTypeHeightMeter                     MeasureType = 4
	MeasureTypeFatFreeMassKilogram             MeasureType = 5
	MeasureTypeFatRatioPercentage              MeasureType = 6
	MeasureTypeFatMassWeightKilogram           MeasureType = 8
	MeasureTypeDiastolicBloodPressuremmHg      MeasureType = 9
	MeasureTypeSystolicBloodPressuremmHg       MeasureType = 10
	MeasureTypeHeartPulseBPM                   MeasureType = 11
	MeasureTypeTemperatureCelsius              MeasureType = 12
	MeasureTypeSPO2                            MeasureType = 54
	MeasureTypeBodyTemperatureCelsius          MeasureType = 71
	MeasureTypeSkinTemperatureCelsius          MeasureType = 73
	MeasureTypeMuscleMassKilogram              MeasureType = 76
	MeasureTypeHydrationKilogram               MeasureType = 77
	MeasureTypeBoneMassKilogram                MeasureType = 88
	MeasureTypePulseWaveVelocityMeterPerSecond MeasureType = 91
	MeasureTypeVo2Max                          MeasureType = 123
	MeasureTypeQRSFromECG                      MeasureType = 135
	MeasureTypePRFromECG                       MeasureType = 136
	MeasureTypeQTFromECG                       MeasureType = 137
	MeasureTypeCorrectedQTFromECG              MeasureType = 138
	MeasureTypeAFibResultFromPPG               MeasureType = 139
)

// Measure is a measure as returned by the Withings API.
type Measure struct {
	Value int64       `json:"value"`
	Type  MeasureType `json:"type"`
	Unit  int         `json:"unit"`
}

// Measures is a slice of Measure structs.
type Measures []Measure

// DecimalValue returns the value of the measure in decimal format by applying the unit value representing the decimal
// location. For example a value of 123 with a unit of 2 would return 1.23.
func (m *Measure) DecimalValue() float64 {
	return float64(m.Value) * math.Pow10(m.Unit)
}

// MeasureGroup is a group of measurements as returned by the Withings API.
type MeasureGroup struct {
	GroupID  int64    `json:"grpid"`
	Attrib   int64    `json:"attrib"`
	Date     int64    `json:"date"`
	Created  int64    `json:"created"`
	Category int64    `json:"category"`
	DeviceID string   `json:"deviceid"`
	Measures Measures `json:"measures"`
	Comment  string   `json:"comment"`
}

// MeasureGroups is a slice of MeasureGroup structs.
type MeasureGroups []MeasureGroup

// GetMeasureResp is the response type returned by the Withings API for a request to obtain measurements.
type GetMeasureResp struct {
	Status   int64          `json:"status"`
	APIError string         `json:"error"`
	Body     GetMeasureBody `json:"body"`
}

// GetMeasureBody is the body of the response returned by the Withings API for a request to obtain measurements.
type GetMeasureBody struct {
	UpdateTime    int64         `json:"updatetime"`
	Timezone      string        `json:"timezone"`
	MeasureGroups MeasureGroups `json:"measuregrps"`
	More          int64         `json:"more"`
	Offset        int64         `json:"offset"`
}

// WeightMeasurement is a parsed withings measurement of the weight type.
type WeightMeasurement struct {
	Pounds    float64
	Kilograms float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToWeight returns a new WeightMeasurement if the measure is of the proper type. If the measure is not a weight
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting WeightMeasurement struct.
func (m *Measure) ToWeight(group *MeasureGroup) *WeightMeasurement {
	if m.Type != MeasureTypeWeightKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &WeightMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.Unix(group.Created, 0)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// Weights returns all the weight values found in every measure group.
func (m MeasureGroups) Weights() []*WeightMeasurement {
	weights := make([]*WeightMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToWeight(&measurementGroup)
			if w != nil {
				weights = append(weights, w)
			}
		}
	}

	return weights
}

// HeightMeasurement is a parsed withings measurement of the height type.
type HeightMeasurement struct {
	Feet     float64
	Meters   float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToHeight returns a new HeightMeasurement if the measure is of the proper type. If the measure is not a Height
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting HeightMeasurement struct.
func (m *Measure) ToHeight(group *MeasureGroup) *HeightMeasurement {
	if m.Type != MeasureTypeHeightMeter {
		return nil
	}

	v := m.DecimalValue()

	w := &HeightMeasurement{
		Feet:   v / 0.3048,
		Meters: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// Heights returns all the Height values found in every measure group.
func (m MeasureGroups) Heights() []*HeightMeasurement {
	heights := make([]*HeightMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToHeight(&measurementGroup)
			if w != nil {
				heights = append(heights, w)
			}
		}
	}

	return heights
}

// FatFreeMassMeasurement is a parsed withings measurement of the fatFreeMass type.
type FatFreeMassMeasurement struct {
	Kilograms float64
	Pounds    float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToFatFreeMass returns a new FatFreeMassMeasurement if the measure is of the proper type. If the measure is not a FatFreeMass
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting FatFreeMassMeasurement struct.
func (m *Measure) ToFatFreeMass(group *MeasureGroup) *FatFreeMassMeasurement {
	if m.Type != MeasureTypeFatFreeMassKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &FatFreeMassMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// FatFreeMasses returns all the FatFreeMass values found in every measure group.
func (m MeasureGroups) FatFreeMasses() []*FatFreeMassMeasurement {
	parsedMeasures := make([]*FatFreeMassMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToFatFreeMass(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// FatRatioMeasurement is a parsed withings measurement of the height type.
type FatRatioMeasurement struct {
	Percentage float64
	Created    time.Time
	DeviceID   string
	GroupID    int64
}

// ToFatRatio returns a new FatRatioMeasurement if the measure is of the proper type. If the measure is not a FatRatio
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting FatRatioMeasurement struct.
func (m *Measure) ToFatRatio(group *MeasureGroup) *FatRatioMeasurement {
	if m.Type != MeasureTypeFatRatioPercentage {
		return nil
	}

	v := m.DecimalValue()

	w := &FatRatioMeasurement{
		Percentage: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// FatRatios returns all the FatRatio values found in every measure group.
func (m MeasureGroups) FatRatios() []*FatRatioMeasurement {
	fatRatios := make([]*FatRatioMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToFatRatio(&measurementGroup)
			if w != nil {
				fatRatios = append(fatRatios, w)
			}
		}
	}

	return fatRatios
}

// FatMassWeightMeasurement is a parsed withings measurement of the height type.
type FatMassWeightMeasurement struct {
	Kilograms float64
	Pounds    float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToFatMassWeight returns a new FatMassWeightMeasurement if the measure is of the proper type. If the measure is not a FatMassWeight
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting FatMassWeightMeasurement struct.
func (m *Measure) ToFatMassWeight(group *MeasureGroup) *FatMassWeightMeasurement {
	if m.Type != MeasureTypeFatMassWeightKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &FatMassWeightMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// FatMassWeights returns all the FatMassWeight values found in every measure group.
func (m MeasureGroups) FatMassWeights() []*FatMassWeightMeasurement {
	fatMassWeights := make([]*FatMassWeightMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToFatMassWeight(&measurementGroup)
			if w != nil {
				fatMassWeights = append(fatMassWeights, w)
			}
		}
	}

	return fatMassWeights
}

// DiastolicBloodPressureMeasurement is a parsed withings measurement of the height type.
type DiastolicBloodPressureMeasurement struct {
	MMHG     float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToDiastolicBloodPressure returns a new DiastolicBloodPressureMeasurement if the measure is of the proper type. If the measure is not a DiastolicBloodPressure
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting DiastolicBloodPressureMeasurement struct.
func (m *Measure) ToDiastolicBloodPressure(group *MeasureGroup) *DiastolicBloodPressureMeasurement {
	if m.Type != MeasureTypeDiastolicBloodPressuremmHg {
		return nil
	}

	v := m.DecimalValue()

	w := &DiastolicBloodPressureMeasurement{
		MMHG: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// DiastolicBloodPressures returns all the DiastolicBloodPressure values found in every measure group.
func (m MeasureGroups) DiastolicBloodPressures() []*DiastolicBloodPressureMeasurement {
	parsedMeasures := make([]*DiastolicBloodPressureMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToDiastolicBloodPressure(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// SystolicBloodPressureMeasurement is a parsed withings measurement of the height type.
type SystolicBloodPressureMeasurement struct {
	MMHG     float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToSystolicBloodPressure returns a new SystolicBloodPressureMeasurement if the measure is of the proper type. If the measure is not a SystolicBloodPressure
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting SystolicBloodPressureMeasurement struct.
func (m *Measure) ToSystolicBloodPressure(group *MeasureGroup) *SystolicBloodPressureMeasurement {
	if m.Type != MeasureTypeSystolicBloodPressuremmHg {
		return nil
	}

	v := m.DecimalValue()

	w := &SystolicBloodPressureMeasurement{
		MMHG: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// SystolicBloodPressures returns all the SystolicBloodPressure values found in every measure group.
func (m MeasureGroups) SystolicBloodPressures() []*SystolicBloodPressureMeasurement {
	parsedMeasures := make([]*SystolicBloodPressureMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToSystolicBloodPressure(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// HeartPulseMeasurement is a parsed withings measurement of the height type.
type HeartPulseMeasurement struct {
	BMP      float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToHeartPulse returns a new HeartPulseMeasurement if the measure is of the proper type. If the measure is not a HeartPulse
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting HeartPulseMeasurement struct.
func (m *Measure) ToHeartPulse(group *MeasureGroup) *HeartPulseMeasurement {
	if m.Type != MeasureTypeHeartPulseBPM {
		return nil
	}

	v := m.DecimalValue()

	w := &HeartPulseMeasurement{
		BMP: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// HeartPulses returns all the HeartPulse values found in every measure group.
func (m MeasureGroups) HeartPulses() []*HeartPulseMeasurement {
	fatRatios := make([]*HeartPulseMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToHeartPulse(&measurementGroup)
			if w != nil {
				fatRatios = append(fatRatios, w)
			}
		}
	}

	return fatRatios
}

// TemperatureMeasurement is a parsed withings measurement of the height type.
type TemperatureMeasurement struct {
	Celsius    float64
	Fahrenheit float64
	Created    time.Time
	DeviceID   string
	GroupID    int64
}

// ToTemperature returns a new TemperatureMeasurement if the measure is of the proper type. If the measure is not a Temperature
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting TemperatureMeasurement struct.
func (m *Measure) ToTemperature(group *MeasureGroup) *TemperatureMeasurement {
	if m.Type != MeasureTypeTemperatureCelsius {
		return nil
	}

	v := m.DecimalValue()

	w := &TemperatureMeasurement{
		Fahrenheit: v*(9.0/5.0) + 32.0,
		Celsius:    v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// Temperatures returns all the Temperature values found in every measure group.
func (m MeasureGroups) Temperatures() []*TemperatureMeasurement {
	parsedMeasures := make([]*TemperatureMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToTemperature(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// SPO2Measurement is a parsed withings measurement of the height type.
type SPO2Measurement struct {
	SPO2     float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToSPO2 returns a new SPO2Measurement if the measure is of the proper type. If the measure is not a SPO2
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting SPO2Measurement struct.
func (m *Measure) ToSPO2(group *MeasureGroup) *SPO2Measurement {
	if m.Type != MeasureTypeSPO2 {
		return nil
	}

	v := m.DecimalValue()

	w := &SPO2Measurement{
		SPO2: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// SPO2s returns all the SPO2 values found in every measure group.
func (m MeasureGroups) SPO2s() []*SPO2Measurement {
	parsedMeasures := make([]*SPO2Measurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToSPO2(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// BodyTemperatureMeasurement is a parsed withings measurement of the height type.
type BodyTemperatureMeasurement struct {
	Celsius    float64
	Fahrenheit float64
	Created    time.Time
	DeviceID   string
	GroupID    int64
}

// ToBodyTemperature returns a new BodyTemperatureMeasurement if the measure is of the proper type. If the measure is not a BodyTemperature
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting BodyTemperatureMeasurement struct.
func (m *Measure) ToBodyTemperature(group *MeasureGroup) *BodyTemperatureMeasurement {
	if m.Type != MeasureTypeBodyTemperatureCelsius {
		return nil
	}

	v := m.DecimalValue()

	w := &BodyTemperatureMeasurement{
		Fahrenheit: v*(9.0/5.0) + 32.0,
		Celsius:    v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// BodyTemperatures returns all the BodyTemperature values found in every measure group.
func (m MeasureGroups) BodyTemperatures() []*BodyTemperatureMeasurement {
	parsedMeasures := make([]*BodyTemperatureMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToBodyTemperature(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// SkinTemperatureMeasurement is a parsed withings measurement of the height type.
type SkinTemperatureMeasurement struct {
	Celsius    float64
	Fahrenheit float64
	Created    time.Time
	DeviceID   string
	GroupID    int64
}

// ToSkinTemperature returns a new SkinTemperatureMeasurement if the measure is of the proper type. If the measure is not a SkinTemperature
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting SkinTemperatureMeasurement struct.
func (m *Measure) ToSkinTemperature(group *MeasureGroup) *SkinTemperatureMeasurement {
	if m.Type != MeasureTypeSkinTemperatureCelsius {
		return nil
	}

	v := m.DecimalValue()

	w := &SkinTemperatureMeasurement{
		Fahrenheit: v*(9.0/5.0) + 32.0,
		Celsius:    v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// SkinTemperatures returns all the SkinTemperature values found in every measure group.
func (m MeasureGroups) SkinTemperatures() []*SkinTemperatureMeasurement {
	parsedMeasures := make([]*SkinTemperatureMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToSkinTemperature(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// MuscleMassMeasurement is a parsed withings measurement of the fatFreeMass type.
type MuscleMassMeasurement struct {
	Kilograms float64
	Pounds    float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToMuscleMass returns a new MuscleMassMeasurement if the measure is of the proper type. If the measure is not a MuscleMass
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting MuscleMassMeasurement struct.
func (m *Measure) ToMuscleMass(group *MeasureGroup) *MuscleMassMeasurement {
	if m.Type != MeasureTypeMuscleMassKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &MuscleMassMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// MuscleMasses returns all the MuscleMass values found in every measure group.
func (m MeasureGroups) MuscleMasses() []*MuscleMassMeasurement {
	parsedMeasures := make([]*MuscleMassMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToMuscleMass(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// HydrationMeasurement is a parsed withings measurement of the fatFreeMass type.
type HydrationMeasurement struct {
	Kilograms float64
	Pounds    float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToHydration returns a new HydrationMeasurement if the measure is of the proper type. If the measure is not a Hydration
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting HydrationMeasurement struct.
func (m *Measure) ToHydration(group *MeasureGroup) *HydrationMeasurement {
	if m.Type != MeasureTypeHydrationKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &HydrationMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// Hydrations returns all the Hydration values found in every measure group.
func (m MeasureGroups) Hydrations() []*HydrationMeasurement {
	parsedMeasures := make([]*HydrationMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToHydration(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// BoneMassMeasurement is a parsed withings measurement of the BoneMassKilogram type.
type BoneMassMeasurement struct {
	Kilograms float64
	Pounds    float64
	Created   time.Time
	DeviceID  string
	GroupID   int64
}

// ToBoneMass returns a new BoneMassMeasurement if the measure is of the proper type. If the measure is not a BoneMass
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting BoneMassMeasurement struct.
func (m *Measure) ToBoneMass(group *MeasureGroup) *BoneMassMeasurement {
	if m.Type != MeasureTypeBoneMassKilogram {
		return nil
	}

	v := m.DecimalValue()

	w := &BoneMassMeasurement{
		Pounds:    v * 2.20462,
		Kilograms: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// BoneMasses returns all the BoneMass values found in every measure group.
func (m MeasureGroups) BoneMasses() []*BoneMassMeasurement {
	parsedMeasures := make([]*BoneMassMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToBoneMass(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// PulseWaveVelocityMeasurement is a parsed withings measurement of the PulseWaveVelocityMeterPerSecond type.
type PulseWaveVelocityMeasurement struct {
	MeterPerSecond float64
	Created        time.Time
	DeviceID       string
	GroupID        int64
}

// ToPulseWaveVelocity returns a new PulseWaveVelocityMeasurement if the measure is of the proper type. If the measure is not a PulseWaveVelocity
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting PulseWaveVelocityMeasurement struct.
func (m *Measure) ToPulseWaveVelocity(group *MeasureGroup) *PulseWaveVelocityMeasurement {
	if m.Type != MeasureTypePulseWaveVelocityMeterPerSecond {
		return nil
	}

	v := m.DecimalValue()

	w := &PulseWaveVelocityMeasurement{
		MeterPerSecond: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// PulseWaveVelocities returns all the PulseWaveVelocity values found in every measure group.
func (m MeasureGroups) PulseWaveVelocities() []*PulseWaveVelocityMeasurement {
	parsedMeasures := make([]*PulseWaveVelocityMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToPulseWaveVelocity(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// Vo2MaxMeasurement is a parsed withings measurement of the Vo2MaxMeterPerSecond type.
type Vo2MaxMeasurement struct {
	Vo2Max   float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToVo2Max returns a new Vo2MaxMeasurement if the measure is of the proper type. If the measure is not a Vo2Max
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting Vo2MaxMeasurement struct.
func (m *Measure) ToVo2Max(group *MeasureGroup) *Vo2MaxMeasurement {
	if m.Type != MeasureTypeVo2Max {
		return nil
	}

	v := m.DecimalValue()

	w := &Vo2MaxMeasurement{
		Vo2Max: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// Vo2Maxes returns all the Vo2Max values found in every measure group.
func (m MeasureGroups) Vo2Maxes() []*Vo2MaxMeasurement {
	parsedMeasures := make([]*Vo2MaxMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToVo2Max(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// QRSMeasurement is a parsed withings measurement of the QRSMeterPerSecond type.
type QRSMeasurement struct {
	QRS      float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToQRS returns a new QRSMeasurement if the measure is of the proper type. If the measure is not a QRS
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting QRSMeasurement struct.
func (m *Measure) ToQRS(group *MeasureGroup) *QRSMeasurement {
	if m.Type != MeasureTypeQRSFromECG {
		return nil
	}

	v := m.DecimalValue()

	w := &QRSMeasurement{
		QRS: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// QRSes returns all the QRS values found in every measure group.
func (m MeasureGroups) QRSes() []*QRSMeasurement {
	parsedMeasures := make([]*QRSMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToQRS(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// PRMeasurement is a parsed withings measurement of the PRMeterPerSecond type.
type PRMeasurement struct {
	PR       float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToPR returns a new PRMeasurement if the measure is of the proper type. If the measure is not a PR
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting PRMeasurement struct.
func (m *Measure) ToPR(group *MeasureGroup) *PRMeasurement {
	if m.Type != MeasureTypePRFromECG {
		return nil
	}

	v := m.DecimalValue()

	w := &PRMeasurement{
		PR: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// PRes returns all the PR values found in every measure group.
func (m MeasureGroups) PRes() []*PRMeasurement {
	parsedMeasures := make([]*PRMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToPR(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// QTMeasurement is a parsed withings measurement of the QTMeterPerSecond type.
type QTMeasurement struct {
	QT       float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToQT returns a new QTMeasurement if the measure is of the proper type. If the measure is not a QT
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting QTMeasurement struct.
func (m *Measure) ToQT(group *MeasureGroup) *QTMeasurement {
	if m.Type != MeasureTypeQTFromECG {
		return nil
	}

	v := m.DecimalValue()

	w := &QTMeasurement{
		QT: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// QTes returns all the QT values found in every measure group.
func (m MeasureGroups) QTes() []*QTMeasurement {
	parsedMeasures := make([]*QTMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToQT(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// CorrectedQTMeasurement is a parsed withings measurement of the CorrectedQTMeterPerSecond type.
type CorrectedQTMeasurement struct {
	QT       float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToCorrectedQT returns a new CorrectedQTMeasurement if the measure is of the proper type. If the measure is not a QT
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting CorrectedQTMeasurement struct.
func (m *Measure) ToCorrectedQT(group *MeasureGroup) *CorrectedQTMeasurement {
	if m.Type != MeasureTypeCorrectedQTFromECG {
		return nil
	}

	v := m.DecimalValue()

	w := &CorrectedQTMeasurement{
		QT: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// CorrectedQTes returns all the QT values found in every measure group.
func (m MeasureGroups) CorrectedQTes() []*CorrectedQTMeasurement {
	parsedMeasures := make([]*CorrectedQTMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToCorrectedQT(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// AfibResultMeasurement is a parsed withings measurement of the AfibResultMeterPerSecond type.
type AfibResultMeasurement struct {
	Value    float64
	Created  time.Time
	DeviceID string
	GroupID  int64
}

// ToAfibResult returns a new AfibResultMeasurement if the measure is of the proper type. If the measure is not a Value
// measurement as defined by the Type field, nil will be returned. If group is non nil then the values
// from the group will be added to the resulting AfibResultMeasurement struct.
func (m *Measure) ToAfibResult(group *MeasureGroup) *AfibResultMeasurement {
	if m.Type != MeasureTypeAFibResultFromPPG {
		return nil
	}

	v := m.DecimalValue()

	w := &AfibResultMeasurement{
		Value: v,
	}

	if group != nil {
		w.Created = time.UnixMilli(group.Created)
		w.DeviceID = group.DeviceID
		w.GroupID = group.GroupID
	}

	return w
}

// AfibResults returns all the Value values found in every measure group.
func (m MeasureGroups) AfibResults() []*AfibResultMeasurement {
	parsedMeasures := make([]*AfibResultMeasurement, 0, 0)

	for _, measurementGroup := range m {
		for _, measurement := range measurementGroup.Measures {
			w := measurement.ToAfibResult(&measurementGroup)
			if w != nil {
				parsedMeasures = append(parsedMeasures, w)
			}
		}
	}

	return parsedMeasures
}

// GetMeasureParam contains the parameters needed to request measurements.
type GetMeasureParam struct {
	// The types of measures to retrieve.
	MeasurementTypes MeasureTypes

	// The category of measurements to retrieve. If not provided MeasureCategoryReal will be used.
	Category MeasureCategory

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
func (p *GetMeasureParam) UpdateQuery(q url.Values) url.Values {
	// Constructing the query parameters based on the param provided.
	q.Set("action", APIActionGetMeasure)
	q.Set("meastypes", p.MeasurementTypes.String())
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

// GetMeasure retrieves measurements for the user represented by the token. Error will be non nil upon an internal
// or api error. If the API returned the error the response will contain the error.
func (c *Client) GetMeasure(ctx context.Context, token AccessToken, param GetMeasureParam) (*GetMeasureResp, error) {

	// Construct authorized request to request data from the API.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIPathGetMeas, nil)
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

	var mResp GetMeasureResp
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
