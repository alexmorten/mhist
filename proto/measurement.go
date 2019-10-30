package proto

import (
	"github.com/alexmorten/mhist/models"
)

// MeasurementFromModel constructs a proto Measurement from the internal measurement definition
func MeasurementFromModel(m models.Measurement) *Measurement {
	pM := &Measurement{}

	if c, ok := m.(*models.Categorical); ok {
		pM.Type = &Measurement_Categorical{Categorical: &Categorical{
			Ts:    c.Ts,
			Value: c.Value,
		}}
	} else if n, ok := m.(*models.Numerical); ok {
		pM.Type = &Measurement_Numerical{Numerical: &Numerical{
			Ts:    n.Ts,
			Value: n.Value,
		}}
	} else {
		return nil
	}

	return pM
}

// ToModel converts the proto measurement to the internal measurement representation
// returns nil if the prot version doesn't contain enough information
func (m *Measurement) ToModel() models.Measurement {
	var modelMeasurent models.Measurement
	if c := m.GetCategorical(); c != nil {
		modelMeasurent = &models.Categorical{
			Ts:    c.Ts,
			Value: c.Value,
		}
	}

	if n := m.GetNumerical(); n != nil {
		modelMeasurent = &models.Numerical{
			Ts:    n.Ts,
			Value: n.Value,
		}
	}

	return modelMeasurent
}
