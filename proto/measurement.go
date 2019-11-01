package proto

import (
	"time"

	"github.com/alexmorten/mhist/models"
)

// MeasurementFromModel constructs a proto Measurement from the internal measurement definition
func MeasurementFromModel(m models.Measurement) *Measurement {
	pM := &Measurement{}

	switch m.(type) {
	case *models.Categorical:
		c := m.(*models.Categorical)
		pM.Type = &Measurement_Categorical{Categorical: &Categorical{
			Ts:    c.Ts,
			Value: c.Value,
		}}
		break
	case *models.Numerical:
		n := m.(*models.Numerical)
		pM.Type = &Measurement_Numerical{Numerical: &Numerical{
			Ts:    n.Ts,
			Value: n.Value,
		}}
		break
	case *models.Raw:
		r := m.(*models.Raw)
		pM.Type = &Measurement_Raw{Raw: &Raw{
			Ts:    r.Ts,
			Value: r.Value,
		}}
		break
	default:
		return nil
	}

	return pM
}

// ToModelWithDefinedTs converts the proto measurement to the internal measurement representation
// returns nil if the protp version doesn't contain enough information
// if the ts is not provided the current time is used
func (m *Measurement) ToModelWithDefinedTs() models.Measurement {
	var modelMeasurent models.Measurement
	if c := m.GetCategorical(); c != nil {
		categorical := &models.Categorical{
			Ts:    c.Ts,
			Value: c.Value,
		}
		if categorical.Ts == 0 {
			categorical.Ts = time.Now().UnixNano()
		}

		modelMeasurent = categorical
	}

	if n := m.GetNumerical(); n != nil {
		numerical := &models.Numerical{
			Ts:    n.Ts,
			Value: n.Value,
		}

		if numerical.Ts == 0 {
			numerical.Ts = time.Now().UnixNano()
		}

		modelMeasurent = numerical
	}

	if r := m.GetRaw(); r != nil {
		raw := &models.Raw{
			Ts:    r.Ts,
			Value: r.Value,
		}

		if raw.Ts == 0 {
			raw.Ts = time.Now().UnixNano()
		}

		modelMeasurent = raw
	}

	return modelMeasurent
}
